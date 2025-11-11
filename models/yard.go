package models

type Yard struct {
	ID     string  `json:"id" gorm:"primaryKey"`
	Name   string  `json:"name"`
	Blocks []Block `json:"blocks" gorm:"foreignKey:YardID"`
}
