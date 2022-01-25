package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teandresmith/injapan-golang-react-backend/helpers"
	"github.com/teandresmith/injapan-golang-react-backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginUser struct{
	FirstName		string		`bson:"first_name" json:"first_name"`
	LastName		string		`bson:"last_name" json:"last_name"`
	Email			string		`bson:"email" json:"email"`
	Token			string		`bson:"token" json:"token"`
	RefreshToken	string		`bson:"refresh_token" json:"refresh_token"`
}

func Login() gin.HandlerFunc{
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

		if admin.Email == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please enter in an email.",
			})
			defer cancel()
			return
		}

		if admin.Password == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please enter in a password.",
			})
			defer cancel()
			return
		}


		var foundAdmin models.Admin
		err := adminCollection.FindOne(ctx, bson.M{"email": admin.Email}).Decode(&foundAdmin)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while querying the admin collection",
				"error": err.Error(),
			})
			return
		}
		

		passwordErr := ValidatePassword(*foundAdmin.Password, *admin.Password)
		if passwordErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Password",
				"error": passwordErr.Error(),
			})
			return
		}

		newToken, refreshToken, err := helpers.GenerateAllToken(*foundAdmin.FirstName, *foundAdmin.LastName, *foundAdmin.Email, true, foundAdmin.AdminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while generating tokens",
				"error": err.Error(),
			})
			return
		}

		updateErr := helpers.UpdateAllTokens(newToken, refreshToken, foundAdmin.AdminID)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while updating tokens",
				"error": err.Error(),
			})
			return
		}

		loginUser := LoginUser{
			FirstName: *foundAdmin.FirstName,
			LastName: *foundAdmin.LastName,
			Email: *foundAdmin.Email,
			Token: newToken,
			RefreshToken: refreshToken,
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login Successful",
			"user": loginUser,
		})

		
	}
}