package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string
	PasswordHash  string
	AgendaSources []AgendaSource `gorm:"constraint:OnDelete:CASCADE;"`
	AgendaItems   []AgendaItem   `gorm:"constraint:OnDelete:CASCADE;"`
	AgendaInvites []AgendaInvite `gorm:"constraint:OnDelete:CASCADE;"`
}
