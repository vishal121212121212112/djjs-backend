package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/gin-gonic/gin"
)

// GetBranchMediaByBranchIDHandler godoc
// @Summary Get Branch Media by Branch ID
// @Description Get all Branch Media records for a specific Branch ID
// @Tags BranchMedia
// @Security ApiKeyAuth
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Param is_child_branch query bool false "Whether this is a child branch (default: false)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/branch-media/branch/{branch_id} [get]
func GetBranchMediaByBranchIDHandler(c *gin.Context) {
	branchIDParam := c.Param("branch_id")
	branchID, err := strconv.ParseUint(branchIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid branch ID"})
		return
	}

	isChildBranch := false
	isChildBranchStr := c.Query("is_child_branch")
	if isChildBranchStr == "true" {
		isChildBranch = true
	}

	mediaList, err := services.GetBranchMediaByBranchID(uint(branchID), isChildBranch)
	// Return empty array if no media found (not an error)
	if err != nil {
		mediaList = []models.BranchMedia{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Branch Media fetched successfully",
		"data":    mediaList,
	})
}

// GetAllBranchMediaHandler retrieves all BranchMedia records
// @Summary Get all Branch Media
// @Description Retrieve all BranchMedia records
// @Tags BranchMedia
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/branch-media [get]
func GetAllBranchMediaHandler(c *gin.Context) {
	medias, err := services.GetAllBranchMedia()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch records"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Branch Media fetched successfully",
		"data":    medias,
	})
}


