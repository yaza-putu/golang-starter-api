package entity

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string    `json:"name" gorm:"name"`
	Email     string    `json:"email" gorm:"email"`
	Password  string    `json:"-" gorm:"password"`
	Avatar    string    `json:"avatar" gorm:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Users []User
