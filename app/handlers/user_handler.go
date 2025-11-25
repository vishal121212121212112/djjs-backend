package handlers

import (
	"net/http"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/gin-gonic/gin"
)

// CreateUserHandler godoc
// @Summary Create a new user
// @Description Create user with auto-generated password (returned in response)
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body models.User true "User payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users [post]
func CreateUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "User created successfully",
		"user":     user,
		"password": user.Password, // show auto-generated password
	})
}

// CreateUserHandler godoc
// @Summary Create a new user
// @Description Create user with auto-generated password (returned in response)
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users [post]

// GetAllUsersHandler godoc
// @Summary     Get all users
// @Tags        Users
// @Security    ApiKeyAuth
// @Produce     json
// @Success     200 {array} models.User
// @Failure     500 {object} map[string]string
// @Router      /api/users [get]
func GetAllUsersHandler(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserSearchHandler godoc
// @Summary     Search users by email or contact number
// @Description Retrieve users based on provided filters (email, contact number, or both).
// @Tags        Users
// @Security    ApiKeyAuth
// @Produce     json
// @Param       email           query string false "User Email"
// @Param       contact_number  query string false "User Contact Number"
// @Success     200 {array} models.User
// @Failure     404 {object} map[string]string
// @Router      /api/users/search [get]
func GetUserSearchHandler(c *gin.Context) {
	email := c.Query("email")
	contact := c.Query("contact_number")

	users, err := services.GetUserSearch(email, contact)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUserHandler godoc
// @Summary Update a user
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body map[string]interface{} true "Updated fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [put]
func UpdateUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateUser(uint(userID), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUserHandler godoc
// @Summary Delete a user (soft delete)
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [delete]
func DeleteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := services.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
