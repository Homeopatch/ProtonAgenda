package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ResourceID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email         string
	PasswordHash  string
	AgendaSources []AgendaSource `gorm:"constraint:OnDelete:CASCADE;"`
	AgendaItems   []AgendaItem   `gorm:"constraint:OnDelete:CASCADE;"`
	AgendaInvites []AgendaInvite `gorm:"constraint:OnDelete:CASCADE;"`
}
