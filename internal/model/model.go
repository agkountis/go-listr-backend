package model

import "github.com/google/uuid"

type List struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name string    `gorm:"not null"`
}

type Item struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Value  string    `gorm:"not null"`
	ListID uuid.UUID `gorm:"not null"`
	List   List
}
