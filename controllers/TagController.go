package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllTags()  gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		results, err := tagCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the tag collection",
				"error": err.Error(),
			})
			return
		}

		var tags []bson.M

		 iterateErr := results.All(ctx, &tags)
		 defer cancel()
		 if iterateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while iterating through tag results",
				"error": iterateErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, tags)
	}
}
