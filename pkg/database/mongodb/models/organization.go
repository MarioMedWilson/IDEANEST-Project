
package models


type OrganizationMember struct {
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	AccessLevel bool `json:"access_level" bson:"access_level"` // true for admin, false for member
}

type Organization struct {
    ID   string `json:"id" bson:"_id"`
    Name string `json:"name" bson:"name"`
		Description string `json:"description" bson:"description"`
		UserID string `json:"user_id" bson:"user_id"`
		OrganizationMembers []OrganizationMember `json:"organization_members" bson:"organization_members"`
}
