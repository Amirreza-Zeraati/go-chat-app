package models

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	UserID uint
	Text   string
}
