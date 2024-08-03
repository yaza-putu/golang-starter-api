package entity

type Token struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RefreshToken struct {
	Token string `json:"refresh_token" form:"refresh_token" validate:"required"`
}
