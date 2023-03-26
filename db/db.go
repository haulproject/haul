package db

import (
	"context"
	"fmt"
	"os"
	"time"

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

	client.Database("haul").Collection(collection).FindOne(ctx, filter).Decode(&result)

	return result, nil

}

// func ReadAll(collection string) ([]*types.Component, error) {
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

	//var components []*types.Component
	var components []*bson.M

	filter := bson.D{primitive.E{}}

	cursor, err := client.Database("haul").Collection(collection).Find(ctx, filter)
	if err != nil {
		return nil, err

	}

	for cursor.Next(ctx) {
		//var component types.Component
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
