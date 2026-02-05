package utils

import "errors"

func ValidateMood(mood int) error {
	if mood < 1 || mood > 5 {
		return errors.New("mood must be between 1 and 5")
	}
	return nil
}
