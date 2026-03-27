package prompt

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/pressctl/cli/internal/utils"
	"github.com/pressctl/cli/pkg/models"
)

// SiteInput holds the input for site creation
type SiteInput struct {
	ServerName    string
	Domain        string
	SiteID        string
	AdminUser     string
	AdminEmail    string
	AdminPassword string
}

// PromptSiteCreate prompts for site creation details
func PromptSiteCreate(servers []models.Server) (*SiteInput, error) {
	input := &SiteInput{}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers available. Add a server first with: press server add")
	}

	provisionedServers := make([]models.Server, 0)
	for _, s := range servers {
		if s.Status == "provisioned" {
			provisionedServers = append(provisionedServers, s)
		}
	}
	if len(provisionedServers) == 0 {
		return nil, fmt.Errorf("no provisioned servers available. Provision a server first with: press server provision <name>")
	}

	// 1. Select server
	serverOpts := make([]huh.Option[int], len(provisionedServers))
	for i, s := range provisionedServers {
		serverOpts[i] = huh.NewOption(
			fmt.Sprintf("%s (%s) — %d sites", s.Name, s.IP, len(s.Sites)), i,
		)
	}
	var serverIndex int
	if err := huh.NewSelect[int]().
		Title("Target server").
		Description("Choose a provisioned server to host this WordPress site").
		Options(serverOpts...).
		Value(&serverIndex).
		Run(); err != nil {
		return nil, normalizeErr(err)
	}
	input.ServerName = provisionedServers[serverIndex].Name

	// 2. Site details form
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Primary domain").
				Description("Main domain for this WordPress site (e.g., example.com)").
				Value(&input.Domain).
				Validate(utils.ValidateDomain),
			huh.NewInput().
				Title("WordPress admin username").
				Value(&input.AdminUser).
				Placeholder("admin"),
			huh.NewInput().
				Title("WordPress admin email").
				Value(&input.AdminEmail).
				Validate(utils.ValidateEmail),
		),
	).Run(); err != nil {
		return nil, normalizeErr(err)
	}

	if strings.TrimSpace(input.AdminUser) == "" {
		input.AdminUser = "admin"
	}
	input.SiteID = generateUniqueSiteID(input.Domain, provisionedServers[serverIndex].Sites)

	// 3. Password
	var useGenerated bool
	if err := huh.NewConfirm().
		Title("Generate a secure password?").
		Affirmative("Yes, generate").
		Negative("I'll enter my own").
		Value(&useGenerated).
		Run(); err != nil {
		return nil, normalizeErr(err)
	}

	if useGenerated {
		input.AdminPassword = GenerateSecurePassword(20)
		fmt.Printf("\n  Generated password: %s\n", input.AdminPassword)
		fmt.Printf("  ⚠️  Save this password — it will not be shown again.\n\n")

		var saved bool
		if err := huh.NewConfirm().
			Title("Have you saved the password?").
			Affirmative("Yes, continue").
			Negative("Not yet").
			Value(&saved).
			Run(); err != nil {
			return nil, normalizeErr(err)
		}
		if !saved {
			return nil, fmt.Errorf("please save the password before continuing")
		}
	} else {
		if err := huh.NewInput().
			Title("WordPress admin password").
			Description("Min 12 chars with uppercase, lowercase, number, and special character").
			EchoMode(huh.EchoModePassword).
			Value(&input.AdminPassword).
			Validate(utils.ValidatePasswordStrength).
			Run(); err != nil {
			return nil, normalizeErr(err)
		}
	}

	// 4. Confirmation
	if err := confirmSiteCreation(input); err != nil {
		return nil, err
	}

	return input, nil
}

func confirmSiteCreation(input *SiteInput) error {
	fmt.Println()
	fmt.Println("  ═══════════════════════════════════")
	fmt.Printf("  Server:      %s\n", input.ServerName)
	fmt.Printf("  Domain:      %s\n", input.Domain)
	fmt.Printf("  Site ID:     %s\n", input.SiteID)
	fmt.Printf("  Admin user:  %s\n", input.AdminUser)
	fmt.Printf("  Admin email: %s\n", input.AdminEmail)
	fmt.Println("  ═══════════════════════════════════")
	fmt.Println()

	var confirm bool
	if err := huh.NewConfirm().
		Title("Create this WordPress site?").
		Affirmative("Yes, create").
		Negative("Cancel").
		Value(&confirm).
		Run(); err != nil {
		return normalizeErr(err)
	}
	if !confirm {
		return fmt.Errorf("site creation cancelled")
	}
	return nil
}

// generateUniqueSiteID creates a unique site ID from the domain, handling collisions
func generateUniqueSiteID(domain string, existingSites []models.Site) string {
	base := generateBaseSiteID(domain)
	candidate := base
	suffix := 2
	for siteIDExists(existingSites, candidate) {
		suffixStr := fmt.Sprintf("%d", suffix)
		maxBaseLen := 16 - len(suffixStr)
		if len(base) > maxBaseLen {
			candidate = base[:maxBaseLen] + suffixStr
		} else {
			candidate = base + suffixStr
		}
		suffix++
	}
	return candidate
}

func generateBaseSiteID(domain string) string {
	tlds := []string{".com", ".net", ".org", ".io", ".co", ".dev", ".app", ".xyz", ".info", ".biz"}
	name := domain
	for _, tld := range tlds {
		name = strings.TrimSuffix(name, tld)
	}
	name = regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(name, "")
	if len(name) > 14 {
		name = name[:14]
	}
	if len(name) < 3 {
		name = "site" + name
	}
	return strings.ToLower(name)
}

func siteIDExists(sites []models.Site, id string) bool {
	for _, site := range sites {
		if site.SiteID == id {
			return true
		}
	}
	return false
}

// GenerateSiteID is exported for use by cmd package in non-interactive mode
func GenerateSiteID(domain string, existingSites []models.Site) string {
	return generateUniqueSiteID(domain, existingSites)
}

// GenerateSecurePassword generates a cryptographically secure random password
func GenerateSecurePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			password[i] = charset[i%len(charset)]
		} else {
			password[i] = charset[num.Int64()]
		}
	}
	return string(password)
}
