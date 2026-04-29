package encryption

import (
	"errors"
	"log"

	"github.com/alexedwards/argon2id"
)

func Hash(text_to_hash string) (string, error) {
	hash, err := argon2id.CreateHash(text_to_hash, argon2id.DefaultParams)
	if err != nil {
		log.Println(err)
		return "", errors.New("Internal error")
	}

	return hash, nil
}

func Match(hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash("minhaSenha", hash)
	if err != nil {
		return false, errors.New("Internal error")
	}

	return match, nil
}
