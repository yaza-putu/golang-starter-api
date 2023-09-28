package entities

import "time"

const AdminId = "eoutKa7q4jdCY"
const UserId = "R8Yg7rKlwmkA4"

type Role struct {
	ID        string `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string `gorrm:"name" json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Roles []Role
