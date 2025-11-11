package models

type Block struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	YardID string `json:"yard_id" gorm:"index"`
	Yard   Yard   `json:"yard" gorm:"foreignKey:YardID"`
	// Kapasitas total block
	TotalSlot int `json:"total_slot"` // Misalnya, 10
	TotalRow  int `json:"total_row"`  // Misalnya, 5
	TotalTier int `json:"total_tier"` // Misalnya, 5
	// Relasi ke rencana
	Plans []YardPlan `json:"plans" gorm:"foreignKey:BlockID"`
	// Occupancy tetap untuk runtime
	Occupancy map[string]bool `json:"-" gorm:"-"` // Key: "slot-row-tier", Value: true jika terisi
}
