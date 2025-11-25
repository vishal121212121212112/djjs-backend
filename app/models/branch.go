package models

import "time"

// swagger:model Branch
type Branch struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"not null" json:"name"`
	Email           string     `gorm:"unique" json:"email,omitempty"`
	CoordinatorName string     `json:"coordinator_name,omitempty"`
	ContactNumber   string     `gorm:"unique;not null" json:"contact_number"`
	EstablishedOn   *time.Time `json:"established_on,omitempty"`
	AashramArea     float64    `json:"aashram_area,omitempty"`
	Country         string     `json:"country,omitempty"`
	State           string     `json:"state,omitempty"`
	District        string     `json:"district,omitempty"`
	City            string     `json:"city,omitempty"`
	Address         string     `json:"address,omitempty"`
	Pincode         string     `json:"pincode,omitempty"`
	PostOffice      string     `json:"post_office,omitempty"`
	PoliceStation   string     `json:"police_station,omitempty"`
	OpenDays        string     `json:"open_days,omitempty"` // e.g. "Mon-Fri"
	DailyStartTime  string     `json:"daily_start_time,omitempty"`
	DailyEndTime    string     `json:"daily_end_time,omitempty"`
	CreatedOn       time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn       *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
	CreatedBy       string     `json:"created_by,omitempty"`
	UpdatedBy       string     `json:"updated_by,omitempty"`
}

// swagger:model BranchInfrastructure
type BranchInfrastructure struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchID  uint       `gorm:"not null" json:"branch_id"`
	Branch    Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Type      string     `gorm:"not null" json:"type"`
	Count     int        `gorm:"not null" json:"count"`
	CreatedOn time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
	CreatedBy string     `json:"created_by,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`
}

func (BranchInfrastructure) TableName() string {
	return "branch_infrastructure"
}

type BranchMember struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberType     string     `gorm:"not null" json:"member_type"`
	Name           string     `gorm:"not null" json:"name"`
	BranchRole     string     `json:"branch_role,omitempty"`
	Responsibility string     `json:"responsibility,omitempty"`
	Age            int        `json:"age,omitempty"`
	DateOfSamarpan *time.Time `json:"date_of_samarpan,omitempty"`
	Qualification  string     `json:"qualification,omitempty"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty"`
	BranchID       uint       `gorm:"not null" json:"branch_id"`
	Branch         Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	CreatedOn      time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn      *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
	CreatedBy      string     `json:"created_by,omitempty"`
	UpdatedBy      string     `json:"updated_by,omitempty"`
}

// Optional: override table name if GORM pluralizes incorrectly
func (BranchMember) TableName() string {
	return "branch_member"
}
