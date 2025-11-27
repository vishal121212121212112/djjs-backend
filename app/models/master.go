package models

type Country struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type State struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	CountryID uint   `json:"country_id"`
}

type City struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `json:"name"`
	StateID uint   `json:"state_id"`
}

type District struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	StateID   uint   `json:"state_id"`
	CountryID uint   `json:"country_id"`
}

type PromotionMaterialType struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	MaterialType string `json:"material_type"`
}
