/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"codeberg.org/haulproject/haul/db"
	"codeberg.org/haulproject/haul/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		// Misc

		e.GET("/v1", handlers.HandleV1)

		e.GET("/v1/healthcheck", handlers.HandleV1Healthcheck)

		// Create

		e.POST("/v1/component", handlers.HandleV1ComponentCreate)

		e.POST("/v1/assembly", handlers.HandleV1AssemblyCreate)

		e.POST("/v1/kit", handlers.HandleV1KitCreate)

		// Read

		e.GET("/v1/component/:component", handlers.HandleV1ComponentRead)

		// List

		e.GET("/v1/component", handlers.HandleV1ComponentList)

		e.GET("/v1/assembly", handlers.HandleV1AssemblyList)

		e.GET("/v1/kit", handlers.HandleV1KitList)

		// Update

		e.PUT("/v1/component/:component", handlers.HandleV1ComponentUpdate)

		// Delete

		e.DELETE("/v1/component/:component", handlers.HandleV1ComponentDelete)

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
