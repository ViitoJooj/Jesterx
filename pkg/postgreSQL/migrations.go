package postgresql

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func Migrations(db *sql.DB) error {
	files, err := filepath.Glob("./migrations/*.sql")
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("erro na migration %s: %w", file, err)
		}
	}

	return nil
}
