package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func DatabaseConnect() *mongo.Client {
	MongoURI := os.Getenv("MONGODB_URI")
	if MongoURI == "" {
		if err := godotenv.Load(); err != nil {
			log.Panic(err)
		}
		MongoURI = os.Getenv("MONGODB_URI")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	return client
	
}

var DatabaseClient *mongo.Client = DatabaseConnect()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("InJapan").Collection(collectionName)

	return collection
}

