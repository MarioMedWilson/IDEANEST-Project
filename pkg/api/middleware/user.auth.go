package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"net/http"
	"strings"
)

func AccessTokenValidate( c *gin.Context ) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	parts := strings.Split(tokenString, " ")
	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		// Add your secret key for token validation here
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
		c.Abort()
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		c.Set("user_id", claims["sub"])
	}
	c.Next()
}
