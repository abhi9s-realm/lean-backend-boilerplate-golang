package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" binding:"required,min=2,max=100"`
	Email     string         `json:"email" binding:"required,email" gorm:"unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
