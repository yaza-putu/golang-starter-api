package unique

import (
	"crypto/rand"
	"unsafe"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var alphabetCharacter = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890$#@+=")

func Uid(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	for i := 0; i < size; i++ {
		b[i] = alphabet[b[i]%byte(len(alphabet))]
	}
	return *(*string)(unsafe.Pointer(&b))
}

func Key(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	for i := 0; i < size; i++ {
		b[i] = alphabetCharacter[b[i]%byte(len(alphabetCharacter))]
	}
	return *(*string)(unsafe.Pointer(&b))
}
