package utils

import (
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail validates an email address
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidPassword validates password strength
func IsValidPassword(password string) bool {
	// At least 8 characters
	if len(password) < 8 {
		return false
	}
	return true
}

// IsValidPhone validates a phone number (basic validation)
func IsValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^[+]?[(]?[0-9]{1,4}[)]?[-\s.]?[(]?[0-9]{1,4}[)]?[-\s.]?[0-9]{1,9}$`)
	return phoneRegex.MatchString(phone)
}
