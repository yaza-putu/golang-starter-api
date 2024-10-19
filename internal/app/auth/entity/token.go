package entity

import "time"

type TokenType string

const ACCESS_TOKEN TokenType = "access_token"
const REFRESH_TOKEN TokenType = "refresh_token"

type Token struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RefreshToken struct {
	DeviceId string `json:"device_id" form:"device_id" validate:"required"`
}

type UserToken struct {
	ID        string `gorm:"primaryKey;type:char(36)" json:"id"`
	DeviceId  string `json:"device_id" gorm:"device_id"`
	UserId    string `json:"user_id" gorm:"user_id"`
	TokenType TokenType
	Token     string    `json:"token" gorm:"token"`
	Device    string    `json:"device" gorm:"device"`
	IP        string    `json:"ip" gorm:"ip"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
