package controllers

import (
	"net/http"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb/models"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/utils"
	"fmt"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb"
)

type UserController struct {
	UserRepository repository.UserRepository
}

func New(userRepository repository.UserRepository) UserController {
	return UserController{
		UserRepository: userRepository,
	}
}


func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingUser, err := uc.UserRepository.GetUserByEmail(c, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(existingUser)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	createdUser, err := uc.UserRepository.CreateUser(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

func (uc *UserController) SignIn(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := uc.UserRepository.GetUserByEmail(c, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(existingUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	
	refreshToken, err := utils.GenerateRefreshToken(existingUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": existingUser,
		"accessToken": accessToken, 
		"refreshToken": refreshToken,
	})
}

func (uc *UserController) UserValidation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User is login successfully!",
	})
}


func (uc *UserController) AccessTokenRefresh (c *gin.Context) {
	var requestBody map[string]string
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	refreshToken, ok := requestBody["refresh_token"]
	utils.RefreshTokenValidate(refreshToken, c)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token not provided"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert user ID to string"})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": newAccessToken})

}

func (uc *UserController) RevokeRefreshToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
	}

	userIDStr, ok := userID.(string)
	if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert user ID to string"})
			return
	}

	redisClient := database.ConnectRedis()
	
	count, err := redisClient.Exists(c, userIDStr).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check refresh token existence"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Refresh token not found"})
		return
	}
	err = redisClient.Del(c, userIDStr).Err()
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke refresh token"})
			return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Refresh token revoked"})
}
