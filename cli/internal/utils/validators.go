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

// reservedSiteIDs are names a site must never take. A site ID becomes the
// site's Linux user and group, MariaDB database and user, PHP-FPM pool and
// Redis account, so it must not name an account or schema the server already
// uses (the playbook also refuses any existing system user or group).
var reservedSiteIDs = map[string]bool{
	// pressctl's own accounts and files
	"pressctl": true, "pressctlbot": true, "admin": true, "default": true,
	"global": true, "www": true, "scheduled": true, "healthcheck": true,
	// MariaDB, Redis and web server
	"mysql": true, "mariadb": true, "sys": true, "informationschema": true,
	"performanceschema": true, "redis": true, "redissites": true, "nginx": true,
	"wwwdata": true, "postfix": true, "debiansysmaint": true,
	// Common Ubuntu system users and groups
	"root": true, "daemon": true, "bin": true, "sync": true, "games": true,
	"man": true, "mail": true, "news": true, "uucp": true, "proxy": true,
	"backup": true, "list": true, "irc": true, "gnats": true, "nobody": true,
	"nogroup": true, "messagebus": true, "syslog": true, "uuidd": true,
	"tcpdump": true, "sshd": true, "landscape": true, "pollinate": true,
	"ubuntu": true, "adm": true, "sudo": true, "staff": true, "users": true,
	"shadow": true, "disk": true, "kmem": true, "utmp": true, "crontab": true,
	"docker": true, "lxd": true, "plugdev": true, "netdev": true, "sslcert": true,
	"polkitd": true, "dnsmasq": true, "clamav": true, "certbot": true,
}

// IsReservedSiteID reports whether name is reserved for the system.
func IsReservedSiteID(name string) bool {
	return reservedSiteIDs[strings.ToLower(name)]
}

// ValidateSiteID validates a site ID (lowercase letters and digits, 3-16
// chars, not a reserved system name). It is used as a Linux user name and a
// MariaDB database name, where '_' and '%' would be grant wildcards.
func ValidateSiteID(name string) error {
	if len(name) < 3 || len(name) > 16 {
		return fmt.Errorf("site ID must be 3-16 characters")
	}
	if !regexp.MustCompile(`^[a-z0-9]+$`).MatchString(name) {
		return fmt.Errorf("site ID must be lowercase letters and digits only")
	}
	if IsReservedSiteID(name) {
		return fmt.Errorf("site ID %q is reserved for the system; choose another", name)
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
