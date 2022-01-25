package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag struct {
	ID    	primitive.ObjectID 		`bson:"_id" json:"_id"`
	Name  	*string      			`bson:"name" json:"name" validate:"required"`
	TagID 	string       			`bson:"tag_id" json:"tag_id"`
}