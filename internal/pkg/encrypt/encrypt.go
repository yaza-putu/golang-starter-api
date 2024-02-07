package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
	"io"
)

func Encrypt(s string) string {
	data := []byte(s)
	key := []byte(config.Key().Passphrase)

	// generate a new aes chiper using 32 byte long key
	c, err := aes.NewCipher(key)
	logger.New(err, logger.SetType(logger.FATAL))

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	logger.New(err, logger.SetType(logger.FATAL))

	// creates a new byte array
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.New(err, logger.SetType(logger.FATAL))
	}
	return string(gcm.Seal(nonce, nonce, data, nil))
}

func Decrypt(e string) string {
	key := []byte(config.Key().Passphrase)

	c, err := aes.NewCipher(key)
	logger.New(err, logger.SetType(logger.FATAL))

	gcm, err := cipher.NewGCM(c)
	logger.New(err, logger.SetType(logger.FATAL))

	nonceSize := gcm.NonceSize()

	if len(e) < nonceSize {
		fmt.Println(err)
	}

	nonce, e := e[:nonceSize], e[nonceSize:]
	plainText, err := gcm.Open(nil, []byte(nonce), []byte(e), nil)
	if err != nil {
		fmt.Println(err)
	}

	return string(plainText)
}
