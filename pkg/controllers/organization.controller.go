package controllers

import (
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb/models"
	"github.com/MarioMedWilson/IDEANEST-Project/pkg/database/mongodb/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrganizationController struct {
	OrganizationRepository repository.OrganizationRepository
}

func NewOrganizationController(organizationRepository repository.OrganizationRepository) OrganizationController {
	return OrganizationController{
		OrganizationRepository: organizationRepository,
	}
}

func (oc *OrganizationController) CreateOrganization(c *gin.Context) {
	var organization models.Organization

	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	createdOrganization, err := oc.OrganizationRepository.CreateOrganization(c, &organization, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdOrganization)
}

func (oc *OrganizationController) GetOrganizations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	organizations, err := oc.OrganizationRepository.GetOrganizations(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, organizations)
}

func (oc *OrganizationController) GetOrganizationID(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	organization, err := oc.OrganizationRepository.GetOrganizationByID(c, c.Param("id"), userID.(string))
	if organization == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found or unauthorized access"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, organization)
}

func (oc *OrganizationController) InviteUser(c *gin.Context) {

	var organizationMember models.OrganizationMember
	if err := c.ShouldBindJSON(&organizationMember); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	
	err := oc.OrganizationRepository.InviteUser(c, c.Param("id"), &organizationMember)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User invited successfully!",

	})
}

func (oc *OrganizationController) OrganizationUserInvitedTo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	organizations, err := oc.OrganizationRepository.OrganizationUserInvitedTo(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, organizations)
}


func (oc *OrganizationController) UpdateOrganization(c *gin.Context) {
	var organization models.Organization
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	organization.ID = c.Param("id")

	result, err := oc.OrganizationRepository.UpdateOrganization(c, userID.(string), &organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Organization updated successfully!",
		"organization": map[string]interface{}{
			"organization_id": result.ID,
			"name": result.Name,
			"description": result.Description,
		},
	})
}

func (oc *OrganizationController) DeleteOrganization(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	err := oc.OrganizationRepository.DeleteOrganization(c, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Organization deleted successfully!",
	})
}

