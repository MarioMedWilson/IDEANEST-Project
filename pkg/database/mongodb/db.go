package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"fmt"
	"os"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var DB *mongo.Database

func ConnectMongoDB(uri, dbName string) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	DB = client.Database(dbName)
	fmt.Println("Connected to MongoDB!")
	return DB
}


func ConnectRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       1,
	})
}
