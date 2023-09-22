package models

import (
  "gorm.io/gorm"
)

type User struct {
  gorm.Model
	Name      string     `gorm:"type:varchar(100);not null"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string     `gorm:"type:varchar(100);not null"`
	Verified  *bool      `gorm:"not null;default:false"`
	Photo     *string    `gorm:"not null;default:'default.png'"`
}

