
package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)


func GenerateAccessToken(userID string) (string, error) {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 10).Unix(), // Access token expires in 24 hours
		"iat": time.Now().Unix(),
		"sub": userID,
	})
	
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

