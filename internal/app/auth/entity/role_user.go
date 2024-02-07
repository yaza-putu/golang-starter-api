package entity

type RoleUser struct {
	ID     string `gorm:"primaryKey;type:char(36)" json:"id"`
	RoleId string `json:"role_id"`
	Role   Role   `json:"-" gorm:"foreignKey:RoleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserId string `json:"user_id"`
	User   User   `json:"-" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
