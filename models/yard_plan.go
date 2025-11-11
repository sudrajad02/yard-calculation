package models

type YardPlan struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	YardID        string  `json:"yard_id" gorm:"index"`
	BlockID       string  `json:"block_id" gorm:"index"`
	PlannedSize   int     `json:"planned_size"`   // 20 atau 40
	PlannedHeight float64 `json:"planned_height"` // 8.6 atau 9.6
	PlannedType   string  `json:"planned_type"`   // DRY, REEFER, dll
	MinSlot       int     `json:"min_slot"`       // Contoh: 4
	MaxSlot       int     `json:"max_slot"`       // Contoh: 7
	MinRow        int     `json:"min_row"`        // Contoh: 1
	MaxRow        int     `json:"max_row"`        // Contoh: 5
	MinTier       int     `json:"min_tier"`       // Contoh: 1
	MaxTier       int     `json:"max_tier"`       // Bergantung kapasitas block
	Yard          Yard    `json:"yard" gorm:"foreignKey:YardID"`
	Block         Block   `json:"block" gorm:"foreignKey:BlockID"`
}
