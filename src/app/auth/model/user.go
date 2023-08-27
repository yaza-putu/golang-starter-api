package model

type User struct {
	ID   string `gorm:"primaryKey;type:char(36)" json:"id"`
	Name string `json:"name" gorm:"name"`
}
