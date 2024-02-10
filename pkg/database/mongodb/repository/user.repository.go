package repository

import (
	"context"

	"errors"
	"golang.org/x/crypto/bcrypt"
	// "example/test/database/models"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	instanceDB *mongo.Database
	ctx context.Context
}

func NewUserRepository(instanceDB *mongo.Database, ctx context.Context) *UserRepository {
	return &UserRepository{
		instanceDB: instanceDB,
		ctx: ctx,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	objectID := primitive.NewObjectID()
	user.ID = objectID.Hex()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	_, err = ur.instanceDB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := ur.instanceDB.Collection("users").FindOne(ctx, bson.M {"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := ur.instanceDB.Collection("users").FindOne(ctx, bson.M {"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}