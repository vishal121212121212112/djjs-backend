package services

import (
	"errors"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// CreateBranchMedia creates a new BranchMedia record
func CreateBranchMedia(media *models.BranchMedia) error {
	return config.DB.Create(media).Error
}

// GetAllBranchMedia retrieves all BranchMedia records
func GetAllBranchMedia() ([]models.BranchMedia, error) {
	var medias []models.BranchMedia
	if err := config.DB.
		Preload("Branch").
		Find(&medias).Error; err != nil {
		return nil, err
	}
	return medias, nil
}

// GetBranchMediaByBranchID retrieves all BranchMedia records by BranchID
func GetBranchMediaByBranchID(branchID uint, isChildBranch bool) ([]models.BranchMedia, error) {
	var mediaList []models.BranchMedia
	if err := config.DB.
		Preload("Branch").
		Where("branch_id = ? AND is_child_branch = ?", branchID, isChildBranch).
		Find(&mediaList).Error; err != nil {
		return nil, errors.New("no branch media found for the given branch ID")
	}
	return mediaList, nil
}

// UpdateBranchMedia updates an existing BranchMedia record
func UpdateBranchMedia(media *models.BranchMedia) error {
	return config.DB.Save(media).Error
}

// DeleteBranchMedia deletes a BranchMedia record
func DeleteBranchMedia(mediaID uint) error {
	return config.DB.Delete(&models.BranchMedia{}, mediaID).Error
}

// GetBranchMediaByID retrieves a BranchMedia record by ID
func GetBranchMediaByID(mediaID uint) (*models.BranchMedia, error) {
	var media models.BranchMedia
	if err := config.DB.First(&media, mediaID).Error; err != nil {
		return nil, errors.New("branch media not found")
	}
	return &media, nil
}


