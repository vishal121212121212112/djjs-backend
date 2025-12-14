package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures user CRUD routes
func SetupUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.POST("", handlers.CreateUserHandler)
		users.GET("", handlers.GetAllUsersHandler)
		users.GET("/:id", handlers.GetUserSearchHandler)
		users.PUT("/:id", handlers.UpdateUserHandler)
		users.DELETE("/:id", handlers.DeleteUserHandler)
		users.POST("/:id/change-password", handlers.ChangePasswordHandler)
		users.POST("/:id/reset-password", handlers.ResetPasswordHandler)
	}
}


