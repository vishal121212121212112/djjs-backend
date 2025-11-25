package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// Create a new event
func CreateEvent(event *models.EventDetails) error {
	event.CreatedOn = time.Now()
	event.UpdatedOn = nil

	if err := config.DB.Create(event).Error; err != nil {
		return err
	}
	return nil
}

// Get all events with type + category
func GetAllEvents() ([]models.EventDetails, error) {
	var events []models.EventDetails

	if err := config.DB.
		Preload("EventType").
		Preload("EventCategory").
		Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

// Search events by type, category, or theme
func SearchEvents(search string) ([]models.EventDetails, error) {
	var events []models.EventDetails

	db := config.DB.Preload("EventType").Preload("EventCategory")

	if search != "" {
		db = db.Where(`
			LOWER(theme) LIKE LOWER(?) OR
			LOWER(scale) LIKE LOWER(?)`,
			"%"+search+"%", "%"+search+"%",
		)
	}

	if err := db.Find(&events).Error; err != nil {
		return nil, errors.New("error fetching events")
	}

	if len(events) == 0 {
		return nil, errors.New("no events found")
	}

	return events, nil
}

// Update event
func UpdateEvent(eventID uint, updatedData map[string]interface{}) error {
	var event models.EventDetails

	if err := config.DB.First(&event, eventID).Error; err != nil {
		return errors.New("event not found")
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&event).Updates(updatedData).Error; err != nil {
		return err
	}

	return nil
}

// Delete event
func DeleteEvent(eventID uint) error {
	if err := config.DB.Delete(&models.EventDetails{}, eventID).Error; err != nil {
		return err
	}
	return nil
}
