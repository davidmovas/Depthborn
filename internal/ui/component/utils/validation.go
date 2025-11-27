package utils

import (
	"fmt"
	"strings"
	"unicode"
)

var (
	// ValidateRequired checks if value is not empty
	ValidateRequired = func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("this field is required")
		}
		return nil
	}

	// ValidateEmail checks if value is valid email format
	ValidateEmail = func(value string) error {
		if !strings.Contains(value, "@") || !strings.Contains(value, ".") {
			return fmt.Errorf("invalid email format")
		}
		return nil
	}

	// ValidateMinLength checks minimum length
	ValidateMinLength = func(min int) func(string) error {
		return func(value string) error {
			if len(value) < min {
				return fmt.Errorf("minimum length is %d characters", min)
			}
			return nil
		}
	}

	// ValidateNumber checks if value is a number
	ValidateNumber = func(value string) error {
		for _, c := range value {
			if !unicode.IsDigit(c) && c != '.' && c != '-' {
				return fmt.Errorf("must be a number")
			}
		}
		return nil
	}
)
