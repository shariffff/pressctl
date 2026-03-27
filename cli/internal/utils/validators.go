package utils

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ValidateDomain validates a domain name format
func ValidateDomain(domain string) error {
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain format (e.g., example.com)")
	}
	return nil
}

// ValidateSiteID validates a site ID (alphanumeric, 3-16 chars)
func ValidateSiteID(name string) error {
	if len(name) < 3 || len(name) > 16 {
		return fmt.Errorf("site ID must be 3-16 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(name) {
		return fmt.Errorf("site ID must be alphanumeric only")
	}
	return nil
}

// ValidateEmail validates an email address format
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePasswordStrength validates password complexity.
// Requires: minimum 12 characters, at least one uppercase, lowercase, number, and special character.
func ValidatePasswordStrength(password string) error {
	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return fmt.Errorf("password must contain at least one %s", strings.Join(missing, ", "))
	}
	return nil
}

// ValidateIP validates an IP address format
func ValidateIP(s string) error {
	if net.ParseIP(s) == nil {
		return fmt.Errorf("invalid IP address format")
	}
	return nil
}

// ValidatePort validates a port number string (1-65535)
func ValidatePort(s string) error {
	port, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("port must be a number")
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}
