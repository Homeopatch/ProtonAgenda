package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AgendaInvite struct {
	ID                uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	UserID            uint
	Description       string
	ExpiresAt         time.Time
	NotBefore         time.Time
	NotAfter          time.Time
	PaddingBefore     time.Duration
	PaddingAfter      time.Duration
	AgendaSources     []AgendaSource     `gorm:"many2many:invite_sources;"`
	ProceduralAgendas []ProceduralAgenda `gorm:"many2many:invite_procedural_agendas;"`
}
