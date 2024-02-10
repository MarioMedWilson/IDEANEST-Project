package repository

import (
	"context"
	// "example/test/database/models"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"errors"
)

type OrganizationRepository struct {
	instanceDB *mongo.Database
	ctx context.Context
}

func NewOrganizationRepository(instanceDB *mongo.Database, ctx context.Context) *OrganizationRepository {
	return &OrganizationRepository{
		instanceDB: instanceDB,
		ctx: ctx,
	}
}

func (or *OrganizationRepository) CreateOrganization(ctx context.Context, organization *models.Organization, userID string) (*models.Organization, error) {
	objectID := primitive.NewObjectID()
	organization.ID = objectID.Hex()

	organization.UserID = userID
	organization.OrganizationMembers = []models.OrganizationMember{} 

	_, err := or.instanceDB.Collection("organizations").InsertOne(ctx, organization)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

func (or *OrganizationRepository) GetOrganizations(ctx context.Context, userID string) ([]models.Organization, error) {
	cursor, err := or.instanceDB.Collection("organizations").Find(ctx, map[string]string{"user_id": userID})
	if err != nil {
		return nil, err
	}
	var organizations []models.Organization
	if err = cursor.All(ctx, &organizations); err != nil {
		return nil, err
	}
	return organizations, nil
}

func (or *OrganizationRepository) GetOrganizationByID(ctx context.Context, id string, userID string) (*models.Organization, error) {
	var organization models.Organization
	err := or.instanceDB.Collection("organizations").FindOne(ctx, map[string]string{"_id": id}).Decode(&organization)

	if organization.UserID != userID {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (or *OrganizationRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := or.instanceDB.Collection("users").FindOne(ctx, bson.M {"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (or *OrganizationRepository) InviteUser(ctx context.Context, organizationID string, organizationMember *models.OrganizationMember) error {
	opts := options.Update().SetUpsert(true)

	User, err := or.GetUserByEmail(ctx, organizationMember.Email)
	if err != nil {
		return err
	}
	if User == nil {
		return fmt.Errorf("User not found")
	}

	organization_members := bson.M{"user_id": User.ID, "name":User.Name, "email": User.Email, "access_level": organizationMember.AccessLevel}
	filter := bson.M{"_id": organizationID}
	update := bson.M{"$push": bson.M{"organization_members": organization_members}}
	result, err := or.instanceDB.Collection("organizations").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

  if result.MatchedCount == 0 {
    return fmt.Errorf("no document found with the provided organization ID")
  }

  if result.ModifiedCount == 0 {
    return fmt.Errorf("the document was not modified, update unsuccessful")
  }

  return nil
}


func (or *OrganizationRepository) UpdateOrganization(ctx context.Context, userID string, organization *models.Organization) (*models.Organization, error) {

	filter := bson.M{"_id": organization.ID}
	
	update := bson.M{"$set": bson.M{
				"name": organization.Name,
				"description": organization.Description,
			},
		}
	_, err := or.instanceDB.Collection("organizations").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	updatedOrg := &models.Organization{}
	err = or.instanceDB.Collection("organizations").FindOne(ctx, filter).Decode(updatedOrg)
	if err != nil {
		return nil, err
	}

	response := &models.Organization{
		ID: updatedOrg.ID,
		Name:           updatedOrg.Name,
		Description:    updatedOrg.Description,
	}

	return response, nil
}

func (or *OrganizationRepository) DeleteOrganization(ctx context.Context, id string) error {
	_, err := or.instanceDB.Collection("organizations").DeleteOne(ctx, map[string]string{"_id": id})
	if err != nil {
		return err
	}
	return nil
}


func (or *OrganizationRepository) OrganizationUserInvitedTo (ctx context.Context, userID string) ([]models.Organization, error) {

	cursor, err := or.instanceDB.Collection("organizations").Find(ctx, map[string]string{"organization_members.user_id": userID})
	if err != nil {
		return nil, err
	}
	var organizations []models.Organization
	if err = cursor.All(ctx, &organizations); err != nil {
		return nil, err
	}
	return organizations, nil
}
