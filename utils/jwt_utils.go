package utils

import (
	"IoTHR-backend/errors"
	"net/http"
	"os"
	"time"

	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var errorInstance = new(errors.ErrorInstance)

type Claims struct {
	UserId primitive.ObjectID `json:"userId"`
	Role   string             `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(userId primitive.ObjectID, role string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserId: userId,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
func ValidateJWT(tokenStr string) (*Claims, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, errorInstance.ReturnError(http.StatusUnauthorized, "Unauthorized")
	}
	if !token.Valid {
		return nil, errorInstance.ReturnError(http.StatusUnauthorized, "Unauthorized")
	}
	return claims, nil
}
