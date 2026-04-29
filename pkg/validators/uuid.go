package validators

import (
	"errors"
	"log"

	utils_repository "github.com/ViitoJooj/Jesterx/pkg/validators/repository"
	"github.com/google/uuid"
)

func Uuid(value string) error {
	id, err := uuid.Parse(value)
	if err != nil {
		log.Println("invalid uuid")
		return errors.New("internal error.")
	}

	if id.Version() != 7 {
		log.Printf("uuid must be v7, got v%d", id.Version())
		return errors.New("internal error.")
	}

	ok, err := utils_repository.Exists(value)
	if err != nil {
		log.Println(err)
		return errors.New("Internal error.")
	}

	if ok {
		log.Println("This uuid already exists")
		return errors.New("Internal error.")
	}

	return nil
}
