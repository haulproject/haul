package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/db"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleV1ComponentTarget(c echo.Context) error {
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

	for key, value := range result {
		if key == "target" {
			return c.JSON(http.StatusOK, value)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"message": fmt.Sprintf("Could not find target for object %s", componentID),
	})
}

func HandleV1ComponentTargetUnset(c echo.Context) error {
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

	empty, err := primitive.ObjectIDFromHex("000000000000000000000000")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprintf("Internal server error"),
		})
	}

	update := bson.D{
		primitive.E{
			Key: "$set", Value: bson.D{
				bson.E{Key: "target", Value: empty}},
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

func HandleV1ComponentTargetSet(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"message": "Not Implemented"})
}
