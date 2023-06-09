// Package handlers implements handlers for the labstack/echo/v4 API routes.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"codeberg.org/haulproject/haul/db"
	"codeberg.org/haulproject/haul/types"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Misc

func HandleV1(c echo.Context) error {
	routes := c.Echo().Routes()
	sort.Slice(routes, func(i, j int) bool { return routes[i].Path < routes[j].Path })
	return c.JSON(http.StatusOK, routes)
}

func HandleV1Healthcheck(c echo.Context) error {
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

// Create

func HandleV1ComponentCreate(c echo.Context) error {
	var components types.Components

	err := c.Bind(&components.Components)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "error",
			"error":   err.Error(),
		})
	}

	result, err := db.CreateComponents(components)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error during db.CreateComponents",
			"error":   err.Error(),
		})
	}

	if len(result.InsertedIDs) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":      "Nothing to do",
			"inserted_ids": result.InsertedIDs,
		})

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Created components",
		"inserted_ids": result.InsertedIDs,
	})
}

func HandleV1AssemblyCreate(c echo.Context) error {
	var assemblies types.Assemblies

	err := c.Bind(&assemblies.Assemblies)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "error",
			"error":   err.Error(),
		})
	}

	result, err := db.CreateAssemblies(assemblies)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error during db.CreateAssemblies",
			"error":   err.Error(),
		})
	}

	if len(result.InsertedIDs) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":      "Nothing to do",
			"inserted_ids": result.InsertedIDs,
		})

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Created assemblies",
		"inserted_ids": result.InsertedIDs,
	})
}

func HandleV1KitCreate(c echo.Context) error {
	var kits types.Kits

	err := c.Bind(&kits.Kits)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "error",
			"error":   err.Error(),
		})
	}

	result, err := db.CreateKits(kits)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error during db.CreateKits",
			"error":   err.Error(),
		})
	}

	if len(result.InsertedIDs) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":      "Nothing to do",
			"inserted_ids": result.InsertedIDs,
		})

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Created kits",
		"inserted_ids": result.InsertedIDs,
	})
}

// Read

func HandleV1ComponentRead(c echo.Context) error {
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

func HandleV1AssemblyRead(c echo.Context) error {
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

	return c.JSON(http.StatusOK, result)
}

func HandleV1KitRead(c echo.Context) error {
	kitID, err := primitive.ObjectIDFromHex(c.Param("kit"))
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

	result, err := db.ReadFromID("kits", kitID)
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

// List

func HandleV1ComponentList(c echo.Context) error {
	components, err := db.ReadAll("components")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, components)
}

func HandleV1AssemblyList(c echo.Context) error {
	assemblies, err := db.ReadAll("assemblies")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, assemblies)
}

func HandleV1KitList(c echo.Context) error {
	kits, err := db.ReadAll("kits")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, kits)
}

// Update

func HandleV1ComponentUpdate(c echo.Context) error {
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

func HandleV1AssemblyUpdate(c echo.Context) error {
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
	var assembly bson.D
	err = bson.Unmarshal(marshalled, &assembly)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}

	validated, err := types.ValidateFields(assembly, types.Assembly{})
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

	result, err := db.UpdateFromID("assemblies", assemblyID, update)
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

func HandleV1KitUpdate(c echo.Context) error {
	kitID, err := primitive.ObjectIDFromHex(c.Param("kit"))
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
	var kit bson.D
	err = bson.Unmarshal(marshalled, &kit)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}

	validated, err := types.ValidateFields(kit, types.Kit{})
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

	result, err := db.UpdateFromID("kits", kitID, update)
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

// Delete

func HandleV1ComponentDelete(c echo.Context) error {
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

func HandleV1AssemblyDelete(c echo.Context) error {
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

	result, err := db.DeleteFromID("assemblies", assemblyID)
	if err != nil {
		// other
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}

func HandleV1KitDelete(c echo.Context) error {
	kitID, err := primitive.ObjectIDFromHex(c.Param("kit"))
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

	result, err := db.DeleteFromID("kits", kitID)
	if err != nil {
		// other
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, result)
}
