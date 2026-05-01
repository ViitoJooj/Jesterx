package validators

import (
	"errors"
	"log"

	utils_repository "github.com/ViitoJooj/Jesterx/pkg/validators/repository"
	"github.com/google/uuid"
)

func Uuid(value uuid.UUID) error {

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
