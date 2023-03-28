/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"codeberg.org/haulproject/haul/db"
	"codeberg.org/haulproject/haul/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the haul server",
	Run: func(cmd *cobra.Command, args []string) {

		mongoUri := viper.GetString("mongo.uri")
		// MongoDB connection

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Use the SetServerAPIOptions() method to set the Stable API version to 1
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

		// Create a new client and connect to the server
		mongoClient, err := mongo.Connect(ctx, opts)

		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err = mongoClient.Disconnect(ctx); err != nil {
				log.Fatal(err)
			}
		}()

		// Ping
		// Send a ping to confirm a successful connection
		log.Println("[info] Trying to ping database...")

		_, err = db.Ping(mongoUri)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("[ok] Database is reachable.")

		// Server

		e := echo.New()

		// Middleware

		e.Pre(middleware.RemoveTrailingSlash())

		// API Routes

		e.GET("/v1", handleV1)

		e.GET("/v1/healthcheck", handleV1Healthcheck)

		e.GET("/v1/component", handleV1ComponentList)

		e.POST("/v1/component", handleV1ComponentCreate)

		e.GET("/v1/component/:component", handleV1ComponentRead)

		e.DELETE("/v1/component/:component", handleV1ComponentDelete)

		e.PUT("/v1/component/:component", handleV1ComponentUpdate)

		// Ready

		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("server.port"))))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// mongo.uri
	serverCmd.Flags().String("mongo-uri", "mongodb://haul:haul@mongo:27017/", "MondoDB connection URI (config: 'mongo.uri')")
	viper.BindPFlag("mongo.uri", serverCmd.Flags().Lookup("mongo-uri"))

	// server.port
	serverCmd.Flags().Int("server-port", 1315, "Server port to expose API (config: 'server.port')")
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("server-port"))
}

// API Handlers

func handleV1(c echo.Context) error {
	routes := c.Echo().Routes()
	sort.Slice(routes, func(i, j int) bool { return routes[i].Path < routes[j].Path })
	return c.JSON(http.StatusOK, routes)
}

func handleV1Healthcheck(c echo.Context) error {
	// Vars

	mongoUri := viper.GetString("mongo.uri")

	// Execution

	// Send a ping to confirm a successful connection
	_, err := db.Ping(mongoUri)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":        "Internal server error",
			"ping_database": "not ok",
		})

	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":        "ok",
		"ping_database": "ok",
	})
}

// POST to insert a component
func handleV1ComponentCreate(c echo.Context) error {
	var component types.Component
	err := c.Bind(&component)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}

	result, err := db.CreateComponent(component)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	message := fmt.Sprintf("Inserted document with _id: %v", result.InsertedID)

	log.Println(message)

	return c.JSON(http.StatusOK, map[string]string{
		"message": message,
	})

}

func handleV1ComponentRead(c echo.Context) error {
	componentID, err := primitive.ObjectIDFromHex(c.Param("component"))
	if err != nil {
		if err == primitive.ErrInvalidHex {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": fmt.Sprintf("%s", err),
			})
		}

		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf("Internal server error"),
		})
	}

	result, err := db.ReadFromID("components", componentID)
	if err != nil || result == nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "No document with specified ObjectID",
			})
		}

		// other
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func handleV1ComponentList(c echo.Context) error {
	components, err := db.ReadAll("components")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, components)
}

func handleV1ComponentDelete(c echo.Context) error {
	componentID, err := primitive.ObjectIDFromHex(c.Param("component"))
	if err != nil {
		if err == primitive.ErrInvalidHex {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": fmt.Sprintf("%s", err),
			})
		}

		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf("Internal server error"),
		})
	}

	result, err := db.DeleteFromID("components", componentID)
	if err != nil {
		// other
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func handleV1ComponentUpdate(c echo.Context) error {
	componentID, err := primitive.ObjectIDFromHex(c.Param("component"))
	if err != nil {
		if err == primitive.ErrInvalidHex {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": fmt.Sprintf("%s", err),
			})
		}

		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf("Internal server error"),
		})
	}

	var data interface{}

	err = c.Bind(&data)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}

	marshalled, err := bson.Marshal(data)

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}
	var component bson.D
	err = bson.Unmarshal(marshalled, &component)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}

	validated, err := types.ValidateFields(component, types.Component{})
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Error during fields validation",
		})
	}

	if validated == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "No valid data to use in update was found, nothing to do",
		})
	}

	update := bson.D{
		primitive.E{
			Key: "$set", Value: validated,
		},
	}

	result, err := db.UpdateFromID("components", componentID, update)
	if err != nil || result == nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "No document with specified ObjectID",
			})
		}

		// other
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	message, err := json.Marshal(result)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error marshalling result",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": string(message)})
}
