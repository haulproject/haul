// Package types implements the types and validation of types necessary for the application to function.
package types

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Component struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Status string   `json:"status"`

	// A component's Target should point to a kit's or assembly's ObjectID
	Target primitive.ObjectID `json:"target"`
}

type Assembly struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Status string   `json:"status"`

	// An assembly's Target should point to a kit's ObjectID
	Target primitive.ObjectID `json:"target"`
}

type Kit struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Status string   `json:"status"`
}

type KitWithID struct {
	ID primitive.ObjectID `json:"_id"`
	Kit
}

/*
Taking a bson.D and a reference interface{}, returns a bson.D that only contains fields whose keys are also valid keys in the reference interface{}.

If an error is encountered, it is returned, and the returned bson.D will be nil.

If no error is encountered, the validated bson.D is returned, and the error will be nil.
*/
func ValidateFields(document bson.D, reference interface{}) (bson.D, error) {
	keys, err := GetFields(reference)
	if err != nil {
		return nil, err
	}

	var validFields bson.D

	for _, value := range document {
		valid := false
		for _, key := range keys {
			if value.Key == key {
				valid = true
			}
		}
		if valid {
			validFields = append(validFields, value)
		}
	}

	return validFields, nil
}

// GetFields returns an unordered list of fields in a reference interface{}
func GetFields(reference interface{}) ([]string, error) {
	var fields []string
	fields_bytes, err := json.Marshal(reference)
	if err != nil {
		return fields, err
	}

	m := make(map[string]interface{})

	if err = json.Unmarshal(fields_bytes, &m); err != nil {
		return fields, err
	}

	for field := range m {
		fields = append(fields, field)
	}

	return fields, nil
}

func GetFieldsOrdered(reference interface{}) ([]string, error) {
	var fields []string

	fields_bytes, err := json.Marshal(reference)
	if err != nil {
		return fields, err
	}

	var fields_map map[string]interface{}

	err = json.Unmarshal(fields_bytes, &fields_map)
	if err != nil {
		return fields, err
	}

	for field := range fields_map {
		fields = append(fields, field)
	}

	return fields, nil
}
