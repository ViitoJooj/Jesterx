package validators

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

func Uuid(value any) error {
	str, ok := value.(string)
	if !ok {
		log.Println("uuid must be a string")
		return errors.New("internal error.")
	}

	id, err := uuid.Parse(str)
	if err != nil {
		log.Println("invalid uuid")
		return errors.New("internal error.")
	}

	if id.Version() != 7 {
		log.Printf("uuid must be v7, got v%d", id.Version())
		return errors.New("internal error.")
	}

	// Precisa validar no banco de dados se já existe o uuid gerado
	// caso já exista, retornar um erro descritivo para recriar o uuid

	return nil
}
