package dependencies

import (
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"gorm.io/gorm"
)

// Dependencies holds all application dependencies
type Dependencies struct {
	DB *gorm.DB
	// Add more dependencies as needed
}

// InitializeDependencies initializes all application dependencies
func InitializeDependencies() (*Dependencies, error) {
	deps := &Dependencies{}

	// Initialize database connection
	config.ConnectDB()
	deps.DB = config.DB

	// Initialize S3 service
	if err := services.InitializeS3(); err != nil {
		return nil, err
	}

	return deps, nil
}

// GetDB returns the database connection
func (d *Dependencies) GetDB() *gorm.DB {
	return d.DB
}

