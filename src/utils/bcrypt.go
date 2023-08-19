package utils

import (
	"github.com/labstack/gommon/log"
	bcrypt2 "golang.org/x/crypto/bcrypt"
)

func Bcrypt(password string) string {
	bytePass := []byte(password)

	hash, err := bcrypt2.GenerateFromPassword(bytePass, bcrypt2.DefaultCost)
	if err != nil {
		log.Error(err)
	}
	return string(hash)
}

func BcryptCheck(password string, hash string) bool {
	bytePass := []byte(password)
	byteHas := []byte(hash)

	err := bcrypt2.CompareHashAndPassword(byteHas, bytePass)

	if err != nil {
		return false
	} else {
		return true
	}
}
