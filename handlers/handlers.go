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

func HandleV1AssemblyCreate(c echo.Context) error {
	var assembly types.Assembly
	err := c.Bind(&assembly)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}

	result, err := db.CreateAssembly(assembly)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	message := fmt.Sprintf("Inserted assembly with _id: %s", result.InsertedID)

	log.Println(message)

	return c.JSON(http.StatusOK, map[string]string{
		"message": message,
	})

}

func HandleV1KitCreate(c echo.Context) error {
	var kit types.Kit
	err := c.Bind(&kit)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Bad request",
		})
	}

	result, err := db.CreateKit(kit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	message := fmt.Sprintf("Inserted document with _id: %s", result.InsertedID)

	log.Println(message)

	return c.JSON(http.StatusOK, map[string]string{
		"message": message,
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

func HandleV1ComponentTags(c echo.Context) error {
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
		if key == "tags" {
			return c.JSON(http.StatusOK, value)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"message": fmt.Sprintf("Could not find tags for object %s", componentID),
	})
}

func HandleV1ComponentTagsClear(c echo.Context) error {
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

	update := bson.D{
		primitive.E{
			Key: "$set", Value: bson.D{
				bson.E{Key: "tags", Value: nil}},
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

// HandleV1ComponentTagsAdd
func HandleV1ComponentTagsAdd(c echo.Context) error {
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

	component, err := db.ReadFromID("components", componentID)
	if err != nil || component == nil {
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

	var tags_old []string
	//tagsFound := false

	for key, value := range component {
		if key == "tags" {
			tags, ok := value.(primitive.A)

			if !ok && value != nil {
				return c.JSON(http.StatusInternalServerError, "[err] Could not iterate over tags")
			}

			//tagsFound = true

			for _, tag := range tags {
				tag_string, ok := tag.(string)
				if !ok {
					return c.JSON(http.StatusInternalServerError, "Cannot cast tag into string")
				}
				tags_old = append(tags_old, tag_string)
			}
		}
	}

	//TODO
	//log.Printf("tags: %s", tags_old)

	/*
		if !tagsFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": fmt.Sprintf("Could not find tags for object %s", componentID),
			})
		}
	*/

	var tags_add bson.A

	err = c.Bind(&tags_add)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	//log.Printf("tags_add: %s", tags_add)

	for _, tag_add := range tags_add {
		present := false
		for _, tag := range tags_old {
			if tag == tag_add {
				present = true
			}

		}

		if !present {
			new_tag, ok := tag_add.(string)
			if !ok {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"message": err.Error(),
				})
			}
			tags_old = append(tags_old, new_tag)
		}
	}

	//log.Printf("tags + tags_add: %s", tags_old)

	// Update
	update := bson.D{
		primitive.E{
			Key: "$set", Value: bson.D{
				bson.E{Key: "tags", Value: tags_old}},
		},
	}

	updateResult, err := db.UpdateFromID("components", componentID, update)
	if err != nil || updateResult == nil {
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

	message, err := json.Marshal(updateResult)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error marshalling updateResult",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": string(message)})
}

// HandleV1ComponentTagsRemove
func HandleV1ComponentTagsRemove(c echo.Context) error {
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

	component, err := db.ReadFromID("components", componentID)
	if err != nil || component == nil {
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

	var tags_old []string
	tagsFound := false

	for key, value := range component {
		if key == "tags" {
			tags, ok := value.(primitive.A)

			if !ok {
				return c.JSON(http.StatusInternalServerError, "[err] Could not iterate over tags")
			}

			tagsFound = true

			for _, tag := range tags {
				tag_string, ok := tag.(string)
				if !ok {
					return c.JSON(http.StatusInternalServerError, "Cannot cast tag into string")
				}
				tags_old = append(tags_old, tag_string)
			}
		}
	}

	if !tagsFound {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": fmt.Sprintf("Could not find tags for object %s", componentID),
		})
	}

	if len(tags_old) == 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("No tags for object %s", componentID),
		})

	}

	var tags_remove bson.A

	err = c.Bind(&tags_remove)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	var tags_new []string

	for _, tag := range tags_old {
		remove := false
		for _, tag_remove := range tags_remove {
			if tag == tag_remove {
				remove = true
			}

		}

		if !remove {
			tags_new = append(tags_new, tag)
		}
	}

	// Update
	update := bson.D{
		primitive.E{
			Key: "$set", Value: bson.D{
				bson.E{Key: "tags", Value: tags_new}},
		},
	}

	updateResult, err := db.UpdateFromID("components", componentID, update)
	if err != nil || updateResult == nil {
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

	message, err := json.Marshal(updateResult)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error marshalling updateResult",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": string(message)})
}
