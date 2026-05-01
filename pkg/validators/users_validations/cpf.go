package validators

import (
	"errors"
	"strings"
	"unicode"
)

func Cpf(cpf string) error {
	cpf = strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, cpf)

	if len(cpf) != 11 {
		return errors.New("Invalid CPF.")
	}

	allEqual := true
	for i := 1; i < 11; i++ {
		if cpf[i] != cpf[0] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return errors.New("Invalid CPF.")
	}

	digits := make([]int, 11)
	for i := 0; i < 11; i++ {
		digits[i] = int(cpf[i] - '0')
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	d1 := (sum * 10) % 11
	if d1 == 10 {
		d1 = 0
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += digits[i] * (11 - i)
	}
	d2 := (sum * 10) % 11
	if d2 == 10 {
		d2 = 0
	}

	if digits[9] != d1 || digits[10] != d2 {
		return errors.New("Invalid CPF.")
	}

	return nil
}
