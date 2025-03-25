package models

import (
	"gorm.io/gorm"
	"net/url"
)

type AgendaSource struct {
	gorm.Model
	Url         url.URL
	UserID      uint
	AgendaItems []AgendaItem `gorm:"constraint:OnDelete:SET NULL;"`
}
