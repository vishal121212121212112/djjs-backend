package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

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
func GetBranchMediaByBranchID(branchID uint) ([]models.BranchMedia, error) {
	var mediaList []models.BranchMedia
	if err := config.DB.
		Preload("Branch").
		Where("branch_id = ?", branchID).
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

// ConvertBranchMediaToPresignedURLs converts BranchMedia items to include presigned URLs
// This function takes a slice of BranchMedia and returns a new slice with presigned URLs
// All media access uses short-lived pre-signed URLs for security
// Items with empty S3Key are skipped with a warning (instead of failing the entire request)
func ConvertBranchMediaToPresignedURLs(ctx context.Context, mediaList []models.BranchMedia) ([]models.BranchMedia, error) {
	result := make([]models.BranchMedia, 0, len(mediaList))
	
	for _, media := range mediaList {
		// Skip items with empty S3Key - log warning but don't fail the entire request
		if media.S3Key == "" {
			log.Printf("WARNING: Skipping branch media item ID %d (branch_id: %d) - empty S3Key. Run backfill migration to populate s3_key from file_url", media.ID, media.BranchID)
			continue
		}
		
		mediaCopy := media
		
		// Generate short-lived presigned URL (15 minutes for gallery listing)
		presignedURL, err := GetPresignedURL(ctx, mediaCopy.S3Key, 15*time.Minute)
		if err != nil {
			// Log error but skip this item instead of failing entire request
			log.Printf("ERROR: Failed to generate presigned URL for branch media ID %d (s3_key: %s): %v", mediaCopy.ID, mediaCopy.S3Key, err)
			continue
		}
		
		// Defensive check: ensure URL is presigned (contains X-Amz-Signature)
		if !strings.Contains(presignedURL, "X-Amz-Signature") && !strings.Contains(presignedURL, "Signature=") {
			log.Printf("ERROR: Generated URL for branch media ID %d does not contain presigned signature: %s", mediaCopy.ID, presignedURL)
			continue
		}
		
		// Store presigned URL in URL field (for JSON serialization)
		// FileURL is internal and not serialized
		mediaCopy.FileURL = presignedURL // Internal storage
		mediaCopy.URL = presignedURL     // JSON response field
		
		result = append(result, mediaCopy)
	}
	
	return result, nil
}


