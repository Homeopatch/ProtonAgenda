package controllers

import (
	"errors"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func ErrorGormToHuma(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return huma.Error404NotFound("Not found", err)
	}
	return err
}
