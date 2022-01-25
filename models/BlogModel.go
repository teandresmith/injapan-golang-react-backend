package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID               	primitive.ObjectID		`bson:"_id" json:"_id"`
	Title            	*string					`bson:"title" json:"title" validate:"required"`
	Subtitle         	*string					`bson:"subtitle" json:"subtitle" validate:"required"`
	Image            	*string					`bson:"image" json:"image" validate:"required"`
	ImageDescription 	*string					`bson:"image_description" json:"image_description" validate:"required"`
	BlogDescription  	*string					`bson:"blog_description" json:"blog_description" validate:"required"`
	Body				*string					`bson:"body" json:"body" validate:"required"`
	CreatedAt   		time.Time				`bson:"created_at" json:"created_at"`
	UpdatedAt			time.Time				`bson:"updated_at" json:"updated_at"`
	Slug				*string					`bson:"slug" json:"slug" validate:"required"`
	BlogID				string					`bson:"blog_id" json:"blog_id"`
	Tags				[]Tag					`bson:"tags" json:"tags" validate:"required"`
}