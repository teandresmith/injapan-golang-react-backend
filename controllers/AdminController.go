package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teandresmith/injapan-golang-react-backend/helpers"
	"github.com/teandresmith/injapan-golang-react-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)



func CreateAdmin() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var admin models.Admin

		if err := c.BindJSON(&admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		if validateErr := validate.Struct(admin); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while validating request body data",
				"error": validateErr.Error(),
			})
			defer cancel()
			return
		}


		admin.ID = primitive.NewObjectID()
		admin.AdminID = admin.ID.Hex()
		admin.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		admin.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		hashPassword, err := HashPassword(*admin.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while validating request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}
		admin.Password = &hashPassword

		token, refreshToken, err := helpers.GenerateAllToken(*admin.FirstName, *admin.LastName, *admin.Email, true, admin.AdminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while generating jwt tokens",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		admin.Token = &token
		admin.RefreshToken = &refreshToken

		insertResults, insertErr := adminCollection.InsertOne(ctx, admin)
		defer cancel()
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while inserting an object in the admin collection",
				"error": insertErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Insertion Successful",
			"result": insertResults,
		})
	}
}

func HashPassword(password string) (string, error) {
	cost := 15
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func ValidatePassword(hashPassword string, providedPassword string) (error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}

func CreateBlog() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var blog models.Blog

		if err := c.BindJSON(&blog); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		if err := validate.Struct(blog); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while validating request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		blog.ID = primitive.NewObjectID()
		blog.BlogID = blog.ID.Hex()
		blog.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		blog.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		results, err := blogCollection.InsertOne(ctx, blog)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while inserting an object in the blog collection",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Insertion Successful",
			"results": results,
		})


	}
}

func FindTagAndGetTag(tagName string) (foundTag models.Tag, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var tag models.Tag

	err := tagCollection.FindOne(ctx, bson.M{"name": tagName}).Decode(&tag)
	defer cancel()
	if err != nil {
		return tag, err.Error()
	}

	return tag, ""
}

func UpdateBlog() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var requestBlog models.Blog

		if err := c.BindJSON(&requestBlog); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		if validateErr := validate.Struct(requestBlog); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while validating the request body data",
				"error": validateErr.Error(),
			})
			defer cancel()
			return
		}

		var updateBlog primitive.D

		if requestBlog.Title != nil {
			updateBlog = append(updateBlog, bson.E{Key: "title", Value: requestBlog.Title})
		}

		if requestBlog.Subtitle != nil {
			updateBlog = append(updateBlog, bson.E{Key: "subtitle", Value: requestBlog.Subtitle})
		}

		if requestBlog.Image != nil {
			updateBlog = append(updateBlog, bson.E{Key: "image", Value: requestBlog.Image})
		}

		if requestBlog.ImageDescription != nil {
			updateBlog = append(updateBlog, bson.E{Key: "image_description", Value: requestBlog.ImageDescription})
		}

		if requestBlog.BlogDescription != nil {
			updateBlog = append(updateBlog, bson.E{Key: "blog_description", Value: requestBlog.BlogDescription})
		}

		if requestBlog.Body != nil {
			updateBlog = append(updateBlog, bson.E{Key: "body", Value: requestBlog.Body})
		}

		if requestBlog.Slug != nil {
			updateBlog = append(updateBlog, bson.E{Key: "slug", Value: requestBlog.Slug})
		}

		if len(requestBlog.Tags) != 0 {
			updateBlog = append(updateBlog, bson.E{Key: "tags", Value: requestBlog.Tags})
		}

		updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateBlog = append(updateBlog, bson.E{Key: "updated_at", Value: updatedAt})

		blogID := c.Param("blogid")

		filter := bson.M{"blog_id": blogID}
		opts := options.Update().SetUpsert(true)
		update := bson.D{{Key: "$set", Value: updateBlog}}

		updateResults, updateErr := blogCollection.UpdateOne(ctx, filter, update, opts)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating an object in the blog collection",
				"error": updateErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Update Successful",
			"results": updateResults,
		})
	}
}

func DeleteBlog() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		blogID := c.Param("blogid")

		results, err := blogCollection.DeleteOne(ctx, bson.M{"blog_id": blogID})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while deleting an object from the blog collection",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Delete Successful",
			"results": results,
		})
	}
}

func CreateTag() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var tag models.Tag

		if err := c.BindJSON(&tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		if err := validate.Struct(tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while validating request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		tag.ID = primitive.NewObjectID()
		tag.TagID = tag.ID.Hex()
		
		results, err := tagCollection.InsertOne(ctx, tag)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while inserting an object in the tag collection",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Insertion Successful",
			"results": results,
		})
	}
}

func UpdateTag() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var tag models.Tag

		if err := c.BindJSON(&tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding the request body data",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		var tagUpdate primitive.D
		tagID := c.Param("tagid")

		if tag.Name != nil {
			tagUpdate = append(tagUpdate, bson.E{Key: "name", Value: tag.Name})
		}


		filter := bson.M{"tag_id": tagID}
		opts := options.Update().SetUpsert(true)
		update := bson.D{{Key: "$set", Value: tagUpdate}}

		updateResult, updateErr := tagCollection.UpdateOne(ctx, filter, update, opts)
		defer cancel()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating an object in the tag collection",
				"error": updateErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Update Successful",
			"result": updateResult,
		})
	}
}

func DeleteTag() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		tagID := c.Param("tagid")

		deleteResult, deleteErr := tagCollection.DeleteOne(ctx, bson.M{"tag_id": tagID})
		defer cancel()
		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while deleting an object in the tag collection",
				"error": deleteErr,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Deletion Successful",
			"result": deleteResult,
		})
	}
}
