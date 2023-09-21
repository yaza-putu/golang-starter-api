package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"yaza/src/config"
)

func Encrypt(s string) string {
	data := []byte(s)
	key := []byte(config.Key().Passphrase)

	// generate a new aes chiper using 32 byte long key
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}

	// creates a new byte array
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	return string(gcm.Seal(nonce, nonce, data, nil))
}

func Decrypt(e string) string {
	key := []byte(config.Key().Passphrase)

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}

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
