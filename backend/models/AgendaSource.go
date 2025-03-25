package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgendaSource struct {
	gorm.Model
	ResourceID  uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Url         string
	UserID      uint
	AgendaItems []AgendaItem `gorm:"constraint:OnDelete:SET NULL;"`
}
