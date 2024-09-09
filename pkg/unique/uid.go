package unique

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	alphabetCharacter = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890$#@+=")
	objectIDCounter   = readRandomUint32()
	processUnique     = processUniqueBytes()
)

// Uid uses an algorithm similar to ObjectId in MongoDB
// to ensure global uniqueness, while also being sorted by timestamp.
func Uid() string {
	var b [12]byte

	// 4 byte timestamp, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))

	// 5 byte process id, big endian
	copy(b[4:9], processUnique[:])

	// 3 byte counter
	putUint24(b[9:12], atomic.AddUint32(&objectIDCounter, 1))

	return hex.EncodeToString(b[:])
}

func Key(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	for i := 0; i < size; i++ {
		b[i] = alphabetCharacter[b[i]%byte(len(alphabetCharacter))]
	}
	return *(*string)(unsafe.Pointer(&b))
}

func processUniqueBytes() [5]byte {
	var b [5]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %w", err))
	}

	return b
}

func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %w", err))
	}

	return (uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
}

func putUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}
