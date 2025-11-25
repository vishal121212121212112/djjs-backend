package models

import "time"

// Volunteer represents volunteer details captured from UI
// swagger:model Volunteer
type Volunteer struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchID      uint       `gorm:"not null" json:"branch_id"`
	Branch        Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	SearchValue   string     `gorm:"column:search_volunteer" json:"search_volunteer,omitempty"`
	VolunteerName string     `gorm:"not null" json:"volunteer_name"`
	NumberOfDays  int        `gorm:"column:number_of_days" json:"number_of_days,omitempty"`
	SevaInvolved  string     `json:"seva_involved,omitempty"`
	MentionSeva   string     `gorm:"column:mention_seva" json:"mention_seva,omitempty"`
	EventID       uint       `json:"event_id"`
	Event         Event      `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
	CreatedOn     time.Time  `json:"created_on,omitempty"`
	UpdatedOn     *time.Time `json:"updated_on,omitempty"`
	CreatedBy     string     `json:"created_by,omitempty"`
	UpdatedBy     string     `json:"updated_by,omitempty"`
}
