package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configures authentication routes
func SetupAuthRoutes(r *gin.RouterGroup) {
	// Public routes
	r.POST("/login", handlers.LoginHandler)

	// Protected routes
	r.POST("/logout", middleware.AuthMiddleware(), handlers.LogoutHandler)
}


