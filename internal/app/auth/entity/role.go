package entity

import "time"

const ADM = "SF%YHm8-XJ^}"
const USR = "Ts6W0l2EU8&v"

type Role struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string    `gorm:"name" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Roles []Role
