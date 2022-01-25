package controllers

import (
	"context"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/teandresmith/injapan-golang-react-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllBlogs() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		results, err := blogCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an issue while querying the blog collection",
				"error": err.Error(),
			})
			return
		}

		var allBlogs []bson.M

		if err := results.All(ctx, &allBlogs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an issue while querying the blog collection",
				"error": err.Error(),
			})
			defer cancel()
			return
		}

		c.JSON(http.StatusOK, allBlogs)
	}
}

func GetBlogByID() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var blog models.Blog
		blogID := c.Param("blogid")

		err := blogCollection.FindOne(ctx, bson.M{"blog_id": blogID}).Decode(&blog)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an issue while querying the blog collection",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, blog)
	}
}

func GetBlogsByParameters() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		sortedBy, exists := c.GetQuery("sortedBy")
		if !exists {
			sortedBy = "created_at"
		}
		
		sortedFlow, exists := c.GetQuery("sortedFlow")
		if !exists {
			sortedFlow = "desc"
		}

		numberOfBlogs, err := strconv.Atoi(c.Query("numberOfBlogs"))
		if err != nil || numberOfBlogs <= 0 {
			numberOfBlogs = 5
		}


		flow := SortedFlow(sortedFlow)

		sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: sortedBy, Value: flow}}}}
		limitStage := bson.D{{Key: "$limit", Value: numberOfBlogs}}

		results, aggregationErr := blogCollection.Aggregate(ctx, mongo.Pipeline{sortStage, limitStage})
		defer cancel()
		if aggregationErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error during aggregation query of the blog collection",
				"error": aggregationErr.Error(),
			})
			return
		}

		var sortedBlogs []bson.M

		iterationErr := results.All(ctx, &sortedBlogs)
		defer cancel()
		if iterationErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while iterating through blog results",
				"error": iterationErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, sortedBlogs)
		
	}
}

type Mail struct {
	Name		string		`bson:"name" json:"name" validate:"required"`
	Email		string		`bson:"email" json:"email" validate:"required"`
	Message		string		`bson:"message" json:"message" validate:"required"`
}

func SendEmail() gin.HandlerFunc{
	return func(c *gin.Context) {
		var mail Mail

		if err := c.BindJSON(&mail); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "There was an error while binding request body data",
				"error": err.Error(),
			})
			return
		}

		if validateErr := validate.Struct(mail); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "THere was an error while validating request body data",
				"error": validateErr.Error(),
			})
			return
		}

		if err := godotenv.Load(); err != nil {
			log.Panic(err)
		}
		username := os.Getenv("EMAIL_USER_NAME")
		password := os.Getenv("EMAIL_PASSWORD")

		auth := smtp.PlainAuth("", username, password, "smtp.gmail.com")
		
		to := []string{username}
		subject := "Subject: InJapan Contact From " + mail.Name + " - " + mail.Email + "\n"
		mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
		emailBody := mail.Message
		msg := []byte(subject + mime + "\n" + emailBody)
		
		err := smtp.SendMail("smtp.gmail.com:587", auth, username, to, msg)
		if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "There was an error while sending the email",
			"error": err.Error(),
		})
		return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Email was sent!",
		})
	}
}

func SortedFlow(sortedFlow string) (int) {
	flow := -1

	if sortedFlow == "asc" {
		flow = 1
	}

	return flow
}
