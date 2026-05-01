package users_validations

import (
	"errors"
	"fmt"
	"strings"
)

func Password(password string) error {
	invalid_chars := []string{"´", "'", "(", ")", "~", "^", "[", "]", "{", "}"}

	if len(password) > 50 {
		return errors.New("Password is too large.")
	}

	if len(password) < 8 {
		return errors.New("Password is too small.")
	}

	for i := 0; i < len(invalid_chars); i++ {
		if strings.Contains(password, invalid_chars[i]) {
			erro := fmt.Sprintf("Password cannot have %s", invalid_chars[i])
			return errors.New(erro)
		}
	}
	return nil
}
