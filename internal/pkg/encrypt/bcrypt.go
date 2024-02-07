package encrypt

import (
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
	bcrypt2 "golang.org/x/crypto/bcrypt"
)

func Bcrypt(password string) string {
	bytePass := []byte(password)

	hash, err := bcrypt2.GenerateFromPassword(bytePass, bcrypt2.DefaultCost)
	logger.New(err, logger.SetType(logger.ERROR))

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
