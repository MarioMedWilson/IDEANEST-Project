package routes

import (
	// "example/test/api/controllers"
	// "example/test/api/middleware"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/controllers"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/api/middleware"
	"github.com/gin-gonic/gin"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)


func SetupUserRoutes(router *gin.RouterGroup, uc controllers.UserController) {
	router.POST("/create", uc.CreateUser)
	router.POST("/signin", uc.SignIn)
	router.POST("/refresh-token", uc.AccessTokenRefresh)
	router.GET("/validate", middleware.AccessTokenValidate, uc.UserValidation)

}

func SetupOrganizationRoutes(router *gin.RouterGroup, oc controllers.OrganizationController) {
	router.POST("/create", middleware.AccessTokenValidate, oc.CreateOrganization)
	router.GET("/", middleware.AccessTokenValidate, oc.GetOrganizations)
	router.GET("/:id", middleware.AccessTokenValidate, oc.GetOrganizationID)
	
	
	router.POST("/:id/invite", middleware.AccessTokenValidate, oc.InviteUser)
	router.GET("/invite", middleware.AccessTokenValidate, oc.OrganizationUserInvitedTo)

	router.PUT("/:id", middleware.AccessTokenValidate, oc.UpdateOrganization)
	router.DELETE("/:id", middleware.AccessTokenValidate, oc.DeleteOrganization)

}

func SetupRoutes (router *gin.RouterGroup, db *mongo.Database, ctx context.Context) {
	router.GET("/test2", controllers.Test)
	router.GET("/fetchall", func(c *gin.Context) {
		controllers.FetchAllData(c, db, ctx)
	})

}
