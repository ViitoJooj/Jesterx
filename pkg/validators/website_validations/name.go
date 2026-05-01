package website_validations

import (
	"errors"
	"strings"
	"unicode"
)

func Name(name string) error {
	if len(name) > 250 {
		return errors.New("Invalid full name.")
	}

	if len(name) < 5 {
		return errors.New("Invalid full name.")
	}

	parts := strings.Fields(name)
	if len(parts) < 2 {
		return errors.New("Invalid full name.")
	}

	for _, p := range parts {
		letterCount := 0
		for i, r := range p {
			switch {
			case unicode.IsLetter(r):
				letterCount++
			case r == '\'' || r == '-':
				if i == 0 || i == len([]rune(p))-1 {
					return errors.New("Invalid full name.")
				}
			default:
				return errors.New("Invalid full name.")
			}
		}

		if letterCount < 2 {
			return errors.New("Invalid full name.")
		}
	}

	return nil
}
