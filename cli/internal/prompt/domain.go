package prompt

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/pressctl/cli/internal/utils"
	"github.com/pressctl/cli/pkg/models"
)

// DomainAddInput holds the input for adding a domain
type DomainAddInput struct {
	ServerName string
	SiteID     string
	Domain     string
	IssueSSL   bool
}

// DomainRemoveInput holds the input for removing a domain
type DomainRemoveInput struct {
	ServerName string
	SiteID     string
	Domain     string
}

// DomainSSLInput holds the input for issuing SSL
type DomainSSLInput struct {
	ServerName string
	SiteID     string
	Domain     string
}

// PromptDomainAdd prompts for domain addition details
func PromptDomainAdd(servers []models.Server) (*DomainAddInput, error) {
	input := &DomainAddInput{}

	type SiteOption struct {
		ServerName string
		Site       models.Site
	}

	var siteOptions []SiteOption
	for _, server := range servers {
		if server.Status == "provisioned" {
			for _, site := range server.Sites {
				siteOptions = append(siteOptions, SiteOption{ServerName: server.Name, Site: site})
			}
		}
	}
	if len(siteOptions) == 0 {
		return nil, fmt.Errorf("no sites available. Create a site first with: press site create")
	}

	opts := make([]huh.Option[int], len(siteOptions))
	for i, opt := range siteOptions {
		opts[i] = huh.NewOption(
			fmt.Sprintf("%s on %s (%d domains)", opt.Site.PrimaryDomain, opt.ServerName, len(opt.Site.Domains)), i,
		)
	}

	var selectedIndex int
	if err := huh.NewSelect[int]().
		Title("Select site").
		Description("Which WordPress site should serve this domain?").
		Options(opts...).
		Value(&selectedIndex).
		Run(); err != nil {
		return nil, normalizeErr(err)
	}
	input.ServerName = siteOptions[selectedIndex].ServerName
	input.SiteID = siteOptions[selectedIndex].Site.SiteID

	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Domain name").
				Description("e.g., www.example.com").
				Value(&input.Domain).
				Validate(func(s string) error {
					if err := utils.ValidateDomain(s); err != nil {
						return err
					}
					for _, domain := range siteOptions[selectedIndex].Site.Domains {
						if domain.Domain == s {
							return fmt.Errorf("domain '%s' already exists on this site", s)
						}
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Issue SSL certificate for this domain?").
				Affirmative("Yes").
				Negative("No, skip SSL").
				Value(&input.IssueSSL),
		),
	).Run(); err != nil {
		return nil, normalizeErr(err)
	}

	return input, nil
}

// PromptDomainRemove prompts for domain removal
func PromptDomainRemove(servers []models.Server) (*DomainRemoveInput, error) {
	input := &DomainRemoveInput{}

	type DomainOption struct {
		ServerName string
		SiteID     string
		Domain     models.Domain
		IsPrimary  bool
	}

	var domainOptions []DomainOption
	for _, server := range servers {
		if server.Status == "provisioned" {
			for _, site := range server.Sites {
				for _, domain := range site.Domains {
					domainOptions = append(domainOptions, DomainOption{
						ServerName: server.Name,
						SiteID:     site.SiteID,
						Domain:     domain,
						IsPrimary:  domain.Domain == site.PrimaryDomain,
					})
				}
			}
		}
	}
	if len(domainOptions) == 0 {
		return nil, fmt.Errorf("no domains available to remove")
	}

	opts := make([]huh.Option[int], len(domainOptions))
	for i, opt := range domainOptions {
		label := fmt.Sprintf("%s — %s on %s", opt.Domain.Domain, opt.SiteID, opt.ServerName)
		if opt.Domain.SSLEnabled {
			label += " [SSL]"
		}
		if opt.IsPrimary {
			label += " (PRIMARY)"
		}
		opts[i] = huh.NewOption(label, i)
	}

	var selectedIndex int
	if err := huh.NewSelect[int]().
		Title("Select domain to remove").
		Options(opts...).
		Value(&selectedIndex).
		Run(); err != nil {
		return nil, normalizeErr(err)
	}

	selected := domainOptions[selectedIndex]

	if selected.IsPrimary {
		fmt.Println()
		fmt.Println("  ⚠️  WARNING: This is the PRIMARY domain for the site.")
		fmt.Println("  Removing it may break the WordPress installation.")
		fmt.Println()

		var confirm bool
		if err := huh.NewConfirm().
			Title("Remove the primary domain?").
			Affirmative("Yes, remove it").
			Negative("Cancel").
			Value(&confirm).
			Run(); err != nil {
			return nil, normalizeErr(err)
		}
		if !confirm {
			return nil, fmt.Errorf("domain removal cancelled")
		}
	}

	input.ServerName = selected.ServerName
	input.SiteID = selected.SiteID
	input.Domain = selected.Domain.Domain

	return input, nil
}

// PromptDomainSSL prompts for SSL certificate issuance
func PromptDomainSSL(servers []models.Server) (*DomainSSLInput, error) {
	input := &DomainSSLInput{}

	type DomainOption struct {
		ServerName string
		SiteID     string
		SiteDomain string
		Domain     models.Domain
	}

	var domainOptions []DomainOption
	for _, server := range servers {
		if server.Status == "provisioned" {
			for _, site := range server.Sites {
				for _, domain := range site.Domains {
					if !domain.SSLEnabled {
						domainOptions = append(domainOptions, DomainOption{
							ServerName: server.Name,
							SiteID:     site.SiteID,
							SiteDomain: site.PrimaryDomain,
							Domain:     domain,
						})
					}
				}
			}
		}
	}
	if len(domainOptions) == 0 {
		return nil, fmt.Errorf("no domains without SSL certificates found")
	}

	opts := make([]huh.Option[int], len(domainOptions))
	for i, opt := range domainOptions {
		opts[i] = huh.NewOption(
			fmt.Sprintf("%s — site: %s on %s", opt.Domain.Domain, opt.SiteDomain, opt.ServerName), i,
		)
	}

	var selectedIndex int
	if err := huh.NewSelect[int]().
		Title("Select domain to issue SSL for").
		Options(opts...).
		Value(&selectedIndex).
		Run(); err != nil {
		return nil, normalizeErr(err)
	}

	selected := domainOptions[selectedIndex]
	input.ServerName = selected.ServerName
	input.SiteID = selected.SiteID
	input.Domain = selected.Domain.Domain

	return input, nil
}
