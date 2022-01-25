package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/teandresmith/injapan-golang-react-backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedTokenDetails struct {
	FirstName 	string 		`bson:"first_name" json:"first_name"`
	LastName  	string 		`bson:"last_name" json:"last_name"`
	Email     	string 		`bson:"email" json:"email"`
	Admin     	bool   		`bson:"admin" json:"admin"`
	Uid			string		`bson:"uid" json:"uid"`
	jwt.RegisteredClaims
}

func Init() string  {
	envVariable := os.Getenv("SECRET_KEY")
	if envVariable == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		envVariable = os.Getenv("SECRET_KEY")
	}
	
	return envVariable
}

var SECRET_KEY = []byte(Init())
var adminCollection *mongo.Collection = database.OpenCollection(database.DatabaseClient, "admins")

func GenerateAllToken(firstName string, lastName string, email string, admin bool, uid string) (signedToken string, signedRefreshToken string, err error) {


	claims := &SignedTokenDetails{
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		Admin: admin,
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour*1)),
			Issuer: "InJapan",
		},
	}

	refreshTokenClaims := &SignedTokenDetails{
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		Admin: admin,
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour*24)),
			Issuer: "InJapan",
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, nil
}

func ValidateAllTokens(tokenString string) (claims *SignedTokenDetails, message string) {

	token, err := jwt.ParseWithClaims(tokenString, &SignedTokenDetails{}, func(token *jwt.Token) (interface{}, error){
		return []byte(SECRET_KEY), nil
	})
	
	if err != nil {
		message = err.Error()
		return &SignedTokenDetails{}, message
	}

	if !token.Valid {
		message = "Token not valid."
		return &SignedTokenDetails{}, message
	}


	claims, ok := token.Claims.(*SignedTokenDetails) 
	if !ok {
		message = err.Error()
		return &SignedTokenDetails{}, message
	}

	
	if !claims.VerifyExpiresAt(time.Now().Local(), true) {
		message = "Token is expired"
		return &SignedTokenDetails{}, message
	}

	return claims, message
}

func UpdateAllTokens(token string, refreshToken string, adminId string) (error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	var updateUserToken primitive.D

	updateUserToken = append(updateUserToken, bson.E{Key: "token", Value: token})
	updateUserToken = append(updateUserToken, bson.E{Key: "refresh_token", Value: refreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateUserToken = append(updateUserToken, bson.E{Key: "updated_at", Value: updated_at})

	filter := bson.M{"admin_id": adminId}
	opt := options.Update().SetUpsert(true)

	_, err := adminCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateUserToken}}, opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return err
	}

	return nil
}