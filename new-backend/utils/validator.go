package utils

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	// At least 8 characters, at least one uppercase, one lowercase, one number
	if len(password) < 8 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasNumber := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}
	return hasUpper && hasLower && hasNumber
}

func SanitizeString(input string) string {
	// Remove leading/trailing whitespace
	input = strings.TrimSpace(input)
	// Remove potential XSS characters
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	return input
}
