package prompt

import (
	"strings"
	"testing"

	"github.com/pressctl/cli/pkg/models"
)

func TestGenerateSiteIDGlobalUniqueness(t *testing.T) {
	servers := []models.Server{
		{Name: "server1", Sites: []models.Site{{SiteID: "sunnyheronsail", PrimaryDomain: "sunnyheronsail.com"}}},
		{Name: "server2", Sites: []models.Site{{SiteID: "greenwavedev", PrimaryDomain: "greenwavedev.com"}}},
	}

	id := GenerateSiteID("sunnyheronsail.com", AllSites(servers))
	if id == "sunnyheronsail" {
		t.Errorf("GenerateSiteID() = %q, want a suffixed ID (collision on another server)", id)
	}
	if id != "sunnyheronsail2" {
		t.Errorf("GenerateSiteID() = %q, want %q", id, "sunnyheronsail2")
	}

	unique := GenerateSiteID("brandnew.com", AllSites(servers))
	if unique != "brandnew" {
		t.Errorf("GenerateSiteID() = %q, want %q", unique, "brandnew")
	}
}

func TestEnsureSiteIDUnique(t *testing.T) {
	servers := []models.Server{
		{Name: "server1", Sites: []models.Site{{SiteID: "sunnyheronsail", PrimaryDomain: "sunnyheronsail.com"}}},
		{Name: "server2", Sites: []models.Site{{SiteID: "greenwavedev", PrimaryDomain: "greenwavedev.com"}}},
	}

	if err := EnsureSiteIDUnique("sunnyheronsail", servers); err == nil {
		t.Error("EnsureSiteIDUnique() = nil, want error for existing site ID on another server")
	} else if !strings.Contains(err.Error(), "server1") {
		t.Errorf("EnsureSiteIDUnique() error = %v, want it to name the colliding server", err)
	}

	if err := EnsureSiteIDUnique("brandnew", servers); err != nil {
		t.Errorf("EnsureSiteIDUnique() = %v, want nil for a fresh site ID", err)
	}
}
