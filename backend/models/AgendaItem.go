package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AgendaItem struct {
	gorm.Model
	ExternalID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	StartTime      time.Time `gorm:"index"`
	EndTime        time.Time
	Description    string
	AgendaSourceID uint
	UserID         uint
}
