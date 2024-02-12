package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb"
	"context"
)

func GenerateRefreshToken(userID string) (string, error) {

	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // Refresh token expires in 7 days
	claims["iat"] = time.Now().Unix()
	claims["sub"] = userID

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	// Save refresh token to Redis
	redisClient := database.ConnectRedis()
	err = redisClient.Set(context.Background(), userID, tokenString, time.Hour*24*7).Err()
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func RefreshTokenValidate(tokenString string, c *gin.Context ) {
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token", "message": err.Error()})
		c.Abort()
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		c.Abort()
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	client := database.ConnectRedis()
	storedToken, err := client.Get(context.Background(), claims["sub"].(string)).Result()
	if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token", "message": err.Error()})
			c.Abort()
			return
	}

	if storedToken != tokenString {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			c.Abort()
			return
	}
	if ok {
		c.Set("user_id", claims["sub"])
	}
	return 
}
