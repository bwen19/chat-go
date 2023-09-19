package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

var (
	isValidName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateName(value string, minLength int, maxLength int) error {
	if err := ValidateString(value, minLength, maxLength); err != nil {
		return err
	}
	if !isValidName(value) {
		return fmt.Errorf("must contain only letters, digits, or underscore")
	}
	return nil
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

func ValidateId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func isOneOf(value string, values []string) error {
	for _, v := range values {
		if value == v {
			return nil
		}
	}
	return fmt.Errorf("must be one of values: %s", strings.Join(values, ""))
}

func ValidateUserRole(value string) error {
	values := []string{"admin", "user"}
	return isOneOf(value, values)
}

func ValidateFriendStatus(value string) error {
	values := []string{"adding", "accepted", "deleted"}
	return isOneOf(value, values)
}

func ValidateRoomCategory(value string) error {
	values := []string{"public", "private", "personal"}
	return isOneOf(value, values)
}

func ValidateRoomRank(value string) error {
	values := []string{"owner", "manager", "member"}
	return isOneOf(value, values)
}

func ValidateMessageKind(value string) error {
	values := []string{"text", "file"}
	return isOneOf(value, values)
}
