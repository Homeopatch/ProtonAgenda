package models

import "gorm.io/gorm"

type ProceduralAgenda struct {
	gorm.Model
	Descriptor  string
	Description string
}
