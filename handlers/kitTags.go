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

func HandleV1KitTags(c echo.Context) error {
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

	for key, value := range result {
		if key == "tags" {
			return c.JSON(http.StatusOK, value)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"message": fmt.Sprintf("Could not find tags for object %s", kitID),
	})
}

func HandleV1KitTagsClear(c echo.Context) error {
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

	update := bson.D{
		primitive.E{
			Key: "$set", Value: bson.D{
				bson.E{Key: "tags", Value: nil}},
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

// HandleV1KitTagsAdd
func HandleV1KitTagsAdd(c echo.Context) error {
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

	kit, err := db.ReadFromID("kits", kitID)
	if err != nil || kit == nil {
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

	for key, value := range kit {
		if key == "tags" {
			tags, ok := value.(primitive.A)

			if !ok && value != nil {
				return c.JSON(http.StatusInternalServerError, "[err] Could not iterate over tags")
			}

			for _, tag := range tags {
				tag_string, ok := tag.(string)
				if !ok {
					return c.JSON(http.StatusInternalServerError, "Cannot cast tag into string")
				}
				tags_old = append(tags_old, tag_string)
			}
		}
	}

	var tags_add bson.A

	err = c.Bind(&tags_add)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

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

	// Update
	update := bson.D{
		primitive.E{
			Key: "$set", Value: bson.D{
				bson.E{Key: "tags", Value: tags_old}},
		},
	}

	updateResult, err := db.UpdateFromID("kits", kitID, update)
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

// HandleV1KitTagsRemove
func HandleV1KitTagsRemove(c echo.Context) error {
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

	kit, err := db.ReadFromID("kits", kitID)
	if err != nil || kit == nil {
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

	for key, value := range kit {
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
			"message": fmt.Sprintf("Could not find tags for object %s", kitID),
		})
	}

	if len(tags_old) == 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("No tags for object %s", kitID),
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

	updateResult, err := db.UpdateFromID("kits", kitID, update)
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
