package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
)

// CreateChildBranchHandler godoc
// @Summary Create a new child branch
// @Description Create a new child branch with all its details
// @Tags Child Branches
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param childBranch body models.ChildBranch true "Child Branch Data"
// @Success 201 {object} models.ChildBranch
// @Failure 400 {object} map[string]string
// @Router /api/child-branches [post]
func CreateChildBranchHandler(c *gin.Context) {
	var childBranch models.ChildBranch

	if err := c.ShouldBindJSON(&childBranch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate parent branch exists
	if childBranch.ParentBranchID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parent_branch_id is required"})
		return
	}

	var parentBranch models.Branch
	if err := config.DB.First(&parentBranch, childBranch.ParentBranchID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent_branch_id"})
		return
	}

	// Always inherit coordinator from parent (override if provided)
	// This ensures coordinator is always the same for child branches
	childBranch.CoordinatorName = parentBranch.CoordinatorName

	if err := services.CreateChildBranch(&childBranch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relations
	createdBranch, err := services.GetChildBranch(childBranch.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created child branch"})
		return
	}

	c.JSON(http.StatusCreated, createdBranch)
}

// GetAllChildBranchesHandler godoc
// @Summary Get all child branches
// @Description Retrieve all child branches with their details
// @Tags Child Branches
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.ChildBranch
// @Router /api/child-branches [get]
func GetAllChildBranchesHandler(c *gin.Context) {
	childBranches, err := services.GetAllChildBranches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, childBranches)
}

// GetChildBranchHandler godoc
// @Summary Get a child branch by ID
// @Description Retrieve a specific child branch by its ID
// @Tags Child Branches
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Child Branch ID"
// @Success 200 {object} models.ChildBranch
// @Failure 404 {object} map[string]string
// @Router /api/child-branches/{id} [get]
func GetChildBranchHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid child branch ID"})
		return
	}

	childBranch, err := services.GetChildBranch(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, childBranch)
}

// GetChildBranchesByParentHandler godoc
// @Summary Get child branches by parent branch ID
// @Description Retrieve all child branches of a specific parent branch
// @Tags Child Branches
// @Security ApiKeyAuth
// @Produce json
// @Param parent_id path int true "Parent Branch ID"
// @Success 200 {array} models.ChildBranch
// @Failure 400 {object} map[string]string
// @Router /api/child-branches/parent/{parent_id} [get]
func GetChildBranchesByParentHandler(c *gin.Context) {
	parentIDParam := c.Param("parent_id")
	parentID, err := strconv.ParseUint(parentIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent branch ID"})
		return
	}

	childBranches, err := services.GetChildBranchesByParent(uint(parentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, childBranches)
}

// UpdateChildBranchHandler godoc
// @Summary Update a child branch
// @Description Update an existing child branch
// @Tags Child Branches
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Child Branch ID"
// @Param childBranch body map[string]interface{} true "Update Data"
// @Success 200 {object} models.ChildBranch
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/child-branches/{id} [put]
func UpdateChildBranchHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid child branch ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the child branch to find its parent
	var childBranch models.ChildBranch
	if err := config.DB.First(&childBranch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "child branch not found"})
		return
	}

	// Get parent branch to inherit coordinator
	var parentBranch models.Branch
	if err := config.DB.First(&parentBranch, childBranch.ParentBranchID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent_branch_id"})
		return
	}

	// Always inherit coordinator from parent (override if provided)
	// This ensures coordinator is always the same for child branches
	updateData["coordinator_name"] = parentBranch.CoordinatorName

	// Remove fields that shouldn't be updated
	delete(updateData, "id")
	delete(updateData, "created_on")
	delete(updateData, "created_by")
	delete(updateData, "parent_branch_id") // Don't allow changing parent

	if err := services.UpdateChildBranch(uint(id), updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated child branch
	updatedBranch, err := services.GetChildBranch(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated child branch"})
		return
	}

	c.JSON(http.StatusOK, updatedBranch)
}

// DeleteChildBranchHandler godoc
// @Summary Delete a child branch
// @Description Delete a child branch by ID
// @Tags Child Branches
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Child Branch ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/child-branches/{id} [delete]
func DeleteChildBranchHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid child branch ID"})
		return
	}

	if err := services.DeleteChildBranch(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "child branch deleted successfully"})
}

// *************************************** Child Branch Infrastructure Handlers ****************************************************** //

// CreateChildBranchInfrastructureHandler godoc
// @Summary Create child branch infrastructure
// @Description Create infrastructure record for a child branch
// @Tags Child Branch Infrastructure
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param infrastructure body models.ChildBranchInfrastructure true "Infrastructure Data"
// @Success 201 {object} models.ChildBranchInfrastructure
// @Failure 400 {object} map[string]string
// @Router /api/child-branches/{id}/infrastructure [post]
func CreateChildBranchInfrastructureHandler(c *gin.Context) {
	var infra models.ChildBranchInfrastructure
	if err := c.ShouldBindJSON(&infra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateChildBranchInfrastructure(&infra); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, infra)
}

// GetChildBranchInfrastructureHandler godoc
// @Summary Get child branch infrastructure
// @Description Get all infrastructure records for a child branch
// @Tags Child Branch Infrastructure
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Child Branch ID"
// @Success 200 {array} models.ChildBranchInfrastructure
// @Router /api/child-branches/{id}/infrastructure [get]
func GetChildBranchInfrastructureHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid child branch ID"})
		return
	}

	infra, err := services.GetInfrastructureByChildBranch(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, infra)
}

// *************************************** Child Branch Member Handlers ****************************************************** //

// CreateChildBranchMemberHandler godoc
// @Summary Create child branch member
// @Description Create a member record for a child branch
// @Tags Child Branch Members
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param member body models.ChildBranchMember true "Member Data"
// @Success 201 {object} models.ChildBranchMember
// @Failure 400 {object} map[string]string
// @Router /api/child-branches/{id}/members [post]
func CreateChildBranchMemberHandler(c *gin.Context) {
	var member models.ChildBranchMember
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateChildBranchMember(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// GetChildBranchMembersHandler godoc
// @Summary Get child branch members
// @Description Get all members of a child branch
// @Tags Child Branch Members
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "Child Branch ID"
// @Success 200 {array} models.ChildBranchMember
// @Router /api/child-branches/{id}/members [get]
func GetChildBranchMembersHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid child branch ID"})
		return
	}

	members, err := services.GetMembersByChildBranch(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

