package handlers

import (
	"fmt"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/db"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleV1AssemblyTarget(c echo.Context) error {
	assemblyID, err := primitive.ObjectIDFromHex(c.Param("assembly"))
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

	result, err := db.ReadFromID("assemblies", assemblyID)
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
		"message": fmt.Sprintf("Could not find target for object %s", assemblyID),
	})
}

func HandleV1AssemblyTargetUnset(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"message": "Not Implemented"})
}

func HandleV1AssemblyTargetSet(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]string{"message": "Not Implemented"})
}
