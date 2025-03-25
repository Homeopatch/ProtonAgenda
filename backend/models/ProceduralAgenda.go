package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProceduralAgenda struct {
	gorm.Model
	ResourceID  uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Descriptor  string
	Description string
}
