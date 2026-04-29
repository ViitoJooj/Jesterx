package utils_repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/supabase"
)

func Exists(value string) (bool, error) {
	if supabase.DB == nil {
		log.Println("database connection is not initialized")
		return false, errors.New("internal error.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	if err := supabase.DB.QueryRowContext(ctx, query, value).Scan(&exists); err != nil {
		log.Printf("error checking uuid existence: %v", err)
		return false, errors.New("internal error.")
	}

	if exists {
		return true, nil
	}

	return false, nil
}
