package validators

import (
	"errors"
	"strings"
)

func Address(country, zipCode, street, number, district, city, state string) error {
	if strings.TrimSpace(country) == "" {
		return errors.New("country is required.")
	}

	if strings.TrimSpace(zipCode) == "" {
		return errors.New("zip code is required.")
	}

	if strings.TrimSpace(street) == "" {
		return errors.New("street is required.")
	}

	if strings.TrimSpace(number) == "" {
		return errors.New("number is required.")
	}

	if strings.TrimSpace(district) == "" {
		return errors.New("district is required.")
	}

	if strings.TrimSpace(city) == "" {
		return errors.New("city is required.")
	}

	if strings.TrimSpace(state) == "" {
		return errors.New("state is required.")
	}

	return nil
}
