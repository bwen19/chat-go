package util

import (
	"fmt"
	"net/mail"
	"regexp"
)

var isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 2, 50); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only letters, digits, or underscore")
	}
	return nil
}

func ValidateName(value string) error {
	return ValidateString(value, 2, 50)
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 50)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}

func ValidateID(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}
