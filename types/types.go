package types

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type Component struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// Taking a bson.D and a reference interface{}, returns a bson.D that only contains fields whose keys are also valid keys in the refenrece interface{}.
// If an error is encountered, it is returned, and the returned bson.D will be nil.
// If no error is encountered, the validated bson.D is returned, and the error will be nil.
func ValidateFields(document bson.D, reference interface{}) (bson.D, error) {
	// Get possible fields
	fields, err := json.Marshal(reference)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})

	if err = json.Unmarshal(fields, &m); err != nil {
		return nil, err
	}

	var keys []string

	for key := range m {
		keys = append(keys, key)
	}

	var validFields bson.D

	for _, value := range document {
		fmt.Printf("%#v\n", value)
		valid := false
		for _, key := range keys {
			if value.Key == key {
				valid = true
			}
		}
		if valid {
			validFields = append(validFields, value)
		} else {
			fmt.Printf("Key [%s] is not a valid key in types.Component\n", value.Key)
		}
	}

	return validFields, nil
}
