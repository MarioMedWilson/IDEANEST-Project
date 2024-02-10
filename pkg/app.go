package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/controllers"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb/repository"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/api/routes"
)

type App struct {
	ucu controllers.UserController
	uru repository.UserRepository

	uco controllers.OrganizationController
	uro repository.OrganizationRepository
	ctx context.Context
	instanceDB *mongo.Database
}

func NewApp() *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoDBURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")

	var instanceDB = database.ConnectMongoDB(mongoDBURI, dbName)
	ctx := context.TODO()

	uru := *repository.NewUserRepository(instanceDB, ctx)
	ucu := controllers.New(uru)
	uro := *repository.NewOrganizationRepository(instanceDB, ctx)
	uco := controllers.NewOrganizationController(uro)

	return &App{
		ucu: ucu,
		uru: uru,
		uco: uco,
		uro: uro,
		ctx: ctx,
		instanceDB: instanceDB,
	}
}

func (app *App) Run() {
	defer app.Close() 
	router := gin.Default()

	router.GET("/test", controllers.Test)
	userRoutes := router.Group("/user")
	routes.SetupUserRoutes(userRoutes, app.ucu)

	organizationRoutes := router.Group("/organization")
	routes.SetupOrganizationRoutes(organizationRoutes, app.uco)

	router.Run()
}

func (app *App) Close() {
	if app.instanceDB != nil {
		err := app.instanceDB.Client().Disconnect(app.ctx)
		if err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		} else {
			log.Println("Disconnected from MongoDB!")
		}
	}
}
