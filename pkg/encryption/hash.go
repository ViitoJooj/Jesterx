package encryption

import (
	"errors"
	"log"

	"github.com/alexedwards/argon2id"
)

func Hash(plainText string) (string, error) {
	hash, err := argon2id.CreateHash(plainText, argon2id.DefaultParams)
	if err != nil {
		log.Println(err)
		return "", errors.New("internal error")
	}

	return hash, nil
}

func Match(plainText, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plainText, hash)
	if err != nil {
		return false, errors.New("internal error")
	}

	return match, nil
}
