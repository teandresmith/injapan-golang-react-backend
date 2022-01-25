package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teandresmith/injapan-golang-react-backend/helpers"
)

func Authentication() gin.HandlerFunc{
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")


		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Token not provided.",
			})
			return
		}


		claims, message := helpers.ValidateAllTokens(token)
		if message != ""{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Token",
				"error": message,
			})
			return
		}

		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("email", claims.Email)
		c.Set("admin", claims.Admin)
		c.Set("uid", claims.Uid)
		
		c.Next()


	}

}
