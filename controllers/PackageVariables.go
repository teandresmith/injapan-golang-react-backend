package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/teandresmith/injapan-golang-react-backend/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var blogCollection *mongo.Collection = database.OpenCollection(database.DatabaseClient, "blogs")
var tagCollection *mongo.Collection = database.OpenCollection(database.DatabaseClient, "tags")
var adminCollection *mongo.Collection = database.OpenCollection(database.DatabaseClient, "admins")
var validate = validator.New()