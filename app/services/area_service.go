package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateArea inserts a new area record
func CreateArea(area *models.Area) error {
	area.DistrictID = uuid.New()
	area.CreatedOn = time.Now()
	area.UpdatedOn = nil

	if err := config.DB.Create(area).Error; err != nil {
		return err
	}
	return nil
}

// GetAllAreas fetches all areas
func GetAllAreas() ([]models.Area, error) {
	var areas []models.Area
	if err := config.DB.Preload("Branch").Find(&areas).Error; err != nil {
		return nil, err
	}
	return areas, nil
}

// GetAreaSearch fetches one area by aresName
func GetAreaSearch(areaName string) ([]models.Area, error) {
	var areas []models.Area
	db := config.DB.Preload("Branch")

	// Apply filters dynamically
	if areaName != "" {
		db = db.Where("LOWER(area_name) LIKE LOWER(?)", "%"+areaName+"%")
	}

	if err := db.Find(&areas).Error; err != nil {
		return nil, errors.New("error fetching areas")
	}

	if len(areas) == 0 {
		return nil, errors.New("no areas found")
	}

	return areas, nil
}

var ErrAreaNotFound = errors.New("area not found")

// UpdateArea updates an area by ID
func UpdateArea(areaID uint, updatedData map[string]interface{}) error {
	var area models.Area
	if err := config.DB.First(&area, areaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAreaNotFound
		}
		return err
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&area).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteArea deletes an area
func DeleteArea(areaID uint) error {
	result := config.DB.Delete(&models.Area{}, areaID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAreaNotFound
	}
	return nil
}
