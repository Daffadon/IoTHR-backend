package utils

import (
	"os"
	"time"

	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
