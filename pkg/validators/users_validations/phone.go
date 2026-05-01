package validators

import (
	"errors"
	"strings"
	"unicode"
)

func Phone(ddi int, ddd int, phone string) error {
	if ddi < 1 || ddi > 999 {
		return errors.New("Invalid phone.")
	}

	phone = strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, phone)

	if len(phone) == 0 {
		return errors.New("Invalid phone.")
	}

	allEqual := true
	for i := 1; i < len(phone); i++ {
		if phone[i] != phone[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return errors.New("Invalid phone.")
	}

	if ddi == 55 {
		if ddd < 11 || ddd > 99 {
			return errors.New("Invalid phone.")
		}

		if len(phone) != 8 && len(phone) != 9 {
			return errors.New("Invalid phone.")
		}

		if len(phone) == 9 && phone[0] != '9' {
			return errors.New("Invalid phone.")
		}

		return nil
	}

	if ddd < 0 || ddd > 9999 {
		return errors.New("Invalid phone.")
	}

	if len(phone) < 4 || len(phone) > 15 {
		return errors.New("Invalid phone.")
	}

	return nil
}
