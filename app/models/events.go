package models

import "time"

type EventType struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
}

type EventCategory struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `json:"name"`
	EventTypeID uint      `json:"event_type_id"`
	EventType   EventType `gorm:"foreignKey:EventTypeID" json:"event_type,omitempty"`
}

type EventDetails struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	EventTypeID uint      `json:"event_type_id"`
	EventType   EventType `gorm:"foreignKey:EventTypeID" json:"event_type,omitempty"`

	EventCategoryID uint          `json:"event_category_id"`
	EventCategory   EventCategory `gorm:"foreignKey:EventCategoryID" json:"event_category,omitempty"`

	Scale           string     `json:"scale,omitempty"`
	Theme           string     `json:"theme,omitempty"`
	StartDate       time.Time  `json:"start_date,omitempty"`
	EndDate         time.Time  `json:"end_date,omitempty"`
	DailyStartTime  *time.Time `json:"daily_start_time,omitempty"`
	DailyEndTime    *time.Time `json:"daily_end_time,omitempty"`
	SpiritualOrator string     `json:"spiritual_orator,omitempty"`

	Country    string `json:"country,omitempty"`
	State      string `json:"state,omitempty"`
	City       string `json:"city,omitempty"`
	District   string `json:"district,omitempty"`
	PostOffice string `json:"post_office,omitempty"`
	Pincode    string `json:"pincode,omitempty"`
	Address    string `json:"address,omitempty"`

	BeneficiaryMen   int `json:"beneficiary_men"`
	BeneficiaryWomen int `json:"beneficiary_women"`
	BeneficiaryChild int `json:"beneficiary_child"`
	InitiationMen    int `json:"initiation_men"`
	InitiationWomen  int `json:"initiation_women"`
	InitiationChild  int `json:"initiation_child"`

	CreatedOn time.Time  `json:"created_on,omitempty"`
	UpdatedOn *time.Time `json:"updated_on,omitempty"`
	CreatedBy string     `json:"created_by,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`
}
