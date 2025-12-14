package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes and groups them together
func SetupRoutes(r *gin.Engine) {
	// Main API group
	api := r.Group("/api")
	{
		// Authentication routes
		SetupAuthRoutes(api)

		// CRUD routes
		SetupAreaRoutes(api)
		SetupUserRoutes(api)
		SetupBranchRoutes(api)
		SetupChildBranchRoutes(api)
		SetupEventRoutes(api)
		SetupPromotionRoutes(api)
		SetupMediaRoutes(api)
		SetupSpecialGuestRoutes(api)
		SetupVolunteerRoutes(api)
		SetupDonationRoutes(api)
		SetupMasterRoutes(api)
		SetupFileRoutes(api)
		SetupBranchMediaRoutes(api)
		SetupChildBranchMediaRoutes(api)
	}
}

