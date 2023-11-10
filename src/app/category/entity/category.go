package entity

import "time"

type Category struct {
	ID        string `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string `gorm:"name" json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Categories []Category
