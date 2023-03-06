package model

import "github.com/google/uuid"

type List struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name string    `gorm:"not null"`
}

type ListItem struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Data   string    `gorm:"not null"`
	ListID uuid.UUID `gorm:"not null"`
	List   List      `json:"-"`
}
