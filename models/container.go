package models

type Container struct {
	ID              uint    `json:"id" gorm:"primaryKey"`
	ContainerNumber string  `json:"container_number" gorm:"uniqueIndex"`
	Size            int     `json:"container_size"`   // 20 atau 40
	Height          float64 `json:"container_height"` // 8.6 atau 9.6
	Type            string  `json:"container_type"`   // DRY, REEFER, OT, dll
	YardID          string  `json:"yard_id"`
	BlockID         string  `json:"block_id"`
	// PlanID untuk mengikat ke rencana tertentu (opsional)
	// YardPlanID        *uint  `json:"yard_plan_id,omitempty"` // Pointer, bisa null
	Slot     int  `json:"slot"`
	Row      int  `json:"row"`
	Tier     int  `json:"tier"`
	IsPlaced bool `json:"isplaced"` // Menandakan apakah kontainer saat ini berada di lapangan
	// Relasi (opsional)
	// YardPlan          *YardPlan `json:"yard_plan,omitempty" gorm:"foreignKey:YardPlanID"`
}
