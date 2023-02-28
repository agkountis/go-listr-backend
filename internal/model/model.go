package model

import "github.com/google/uuid"

type List struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name string
}

type Item struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	ListID uuid.UUID
	List   List
}
