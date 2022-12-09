package util

import (
	"crypto/rand"
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"unsafe"

	"github.com/matthewhartstonge/argon2"
)

func Hash(passwd []byte) ([]byte, error) {
	argon := argon2.DefaultConfig()

	buf, err := argon.HashEncoded(passwd)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func HashString(passwd []byte) (string, error) {
	buf, err := Hash(passwd)
	if err != nil {
		return "", err
	}

	return *(*string)(unsafe.Pointer(&buf)), nil
}

func VerifyEncoded(passwd []byte, encoded []byte) bool {
	ok, err := argon2.VerifyEncoded(passwd, encoded)
	if err != nil {
		return false
	}

	return ok
}

func GenerateKey() string {
	buff := make([]byte, 32)
	if _, err := rand.Read(buff); err != nil {
		jww.FATAL.Fatal(err)
	}

	return fmt.Sprintf("%x", buff)
}
