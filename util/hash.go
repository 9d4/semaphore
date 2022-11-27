package util

import (
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
