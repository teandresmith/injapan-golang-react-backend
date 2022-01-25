package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID         		primitive.ObjectID  	`bson:"_id"`
	FirstName  		*string 				`bson:"first_name" json:"first_name" validate:"required"`
	LastName   		*string 				`bson:"last_name" json:"last_name" validate:"required"`
	Email      		*string 				`bson:"email" json:"email" validate:"email,required"`
	Password   		*string 				`bson:"Password" json:"Password" validate:"required"`
	Token			*string 				`bson:"token" json:"token"`
	RefreshToken	*string					`bson:"refresh_token" json:"refresh_token"`
	Created_At 		time.Time 				`bson:"created_at" json:"created_at"`
	Updated_At 		time.Time 				`bson:"updated_at" json:"updated_at"`
	AdminID    		string 					`bson:"admin_id" json:"admin_id"`
}