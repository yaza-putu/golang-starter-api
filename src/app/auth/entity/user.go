package entity

import (
	"time"
)

type User struct {
	ID        string `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string `json:"name" gorm:"name"`
	Email     string `json:"email" gorm:"email"`
	Password  string `json:"password" gorm:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Users []User
