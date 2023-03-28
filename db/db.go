package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Ping(uri string) (bson.M, error) {
	// MongoDB connection

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(viper.GetString("mongo.uri")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}()

	var result bson.M

	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func CreateComponent(component types.Component) (*mongo.InsertOneResult, error) {
	mongoUri := viper.GetString("mongo.uri")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}()
	coll := client.Database("haul").Collection("components")

	/*
		var component types.Component
		err = c.Bind(&component)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Bad request",
			})
		}
	*/

	//TODO Add JSON validation

	if component.Name == "" {
		return nil, errors.New("component.Name cannot be empty")
	}

	result, err := coll.InsertOne(context.TODO(), component)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ReadFromID(collection string, id primitive.ObjectID) (bson.M, error) {

	// MongoDB connection

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(viper.GetString("mongo.uri")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}()

	var result bson.M

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	err = client.Database("haul").Collection(collection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil

}

func ReadAll(collection string) ([]*bson.M, error) {
	// MongoDB connection

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(viper.GetString("mongo.uri")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}()

	var components []*bson.M

	filter := bson.D{primitive.E{}}

	cursor, err := client.Database("haul").Collection(collection).Find(ctx, filter)
	if err != nil {
		return nil, err

	}

	for cursor.Next(ctx) {
		var component bson.M
		err := cursor.Decode(&component)
		if err != nil {
			return nil, err
		}
		components = append(components, &component)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(ctx)

	return components, nil
}

func DeleteFromID(collection string, id primitive.ObjectID) (*mongo.DeleteResult, error) {

	// MongoDB connection

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(viper.GetString("mongo.uri")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}()

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	result, err := client.Database("haul").Collection(collection).DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func UpdateFromID(collection string, id primitive.ObjectID, data bson.D) (*mongo.UpdateResult, error) {
	// Empty name validation

	dataBytes, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}

	var component map[string]map[string]interface{}

	err = bson.Unmarshal(dataBytes, &component)
	if err != nil {
		return nil, err
	}

	// Validation for trying to empty the name
	for _, element := range component {
		for key, value := range element {
			if key == "name" && value == "" {
				return nil, errors.New("Cannot insert empty string into Component.Name")
			}
		}
	}

	// MongoDB connection

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(viper.GetString("mongo.uri")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}()

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	result, err := client.Database("haul").Collection(collection).UpdateOne(ctx, filter, data)
	if err != nil {
		return nil, err
	}

	return result, nil

}
