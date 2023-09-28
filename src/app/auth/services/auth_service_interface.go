package services

import "yaza/src/http/response"

type AuthInterface interface {
	Login(email string, password string) response.DataApi
	Refresh(oToken string) response.DataApi
}
