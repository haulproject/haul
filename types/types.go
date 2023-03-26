package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Component struct {
	ID   primitive.ObjectID `mongo:"_id" json:"_id"`
	Name string             `json:"name"`
	Tags []string           `json:"tags"`
}
