package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	// Version information
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"

	// Global flags
	Verbose bool
	DryRun  bool
)

const repoURL = "github.com/shariffff/pressctl"

// menuModel is a Bubble Tea model for the main menu.
// Supports arrow keys, vim keys (j/k), number shortcuts, and live type-to-filter.
type menuModel struct {
	items   []string
	visible []int  // indices into items that match current filter
	cursor  int    // position within visible
	filter  string // current search string
	chosen  int    // -1 = pending, -2 = cancelled, >=0 = index into items
}

func newMenuModel(items []string) menuModel {
	m := menuModel{items: items, chosen: -1}
	m.visible = allIndices(items)
	return m
}

func allIndices(items []string) []int {
	idx := make([]int, len(items))
	for i := range items {
		idx[i] = i
	}
	return idx
}

func filterIndices(items []string, query string) []int {
	if query == "" {
		return allIndices(items)
	}
	lower := strings.ToLower(query)
	var result []int
	for i, item := range items {
		if strings.Contains(strings.ToLower(item), lower) {
			result = append(result, i)
		}
	}
	return result
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.chosen = -2
			return m, tea.Quit

		case "esc":
			if m.filter != "" {
				// Clear filter first, don't exit
				m.filter = ""
				m.visible = allIndices(m.items)
				m.cursor = 0
			} else {
				m.chosen = -2
				return m, tea.Quit
			}

		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.visible = filterIndices(m.items, m.filter)
				if m.cursor >= len(m.visible) {
					m.cursor = max(0, len(m.visible)-1)
				}
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.visible)-1 {
				m.cursor++
			}

		case "enter":
			if len(m.visible) > 0 {
				m.chosen = m.visible[m.cursor]
				return m, tea.Quit
			}

		default:
			s := msg.String()

			// Number shortcut: 1–9 selects nth visible item
			if len(s) == 1 && s >= "1" && s <= "9" {
				n, _ := strconv.Atoi(s)
				if n >= 1 && n <= len(m.visible) {
					m.chosen = m.visible[n-1]
					return m, tea.Quit
				}
			}

			// Any other printable single character: append to filter
			if len(s) == 1 && s >= " " {
				m.filter += s
				m.visible = filterIndices(m.items, m.filter)
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	if m.chosen != -1 {
		return ""
	}

	const (
		cyan   = "\033[36m"
		bold   = "\033[1m"
		dim    = "\033[2m"
		yellow = "\033[33m"
		reset  = "\033[0m"
	)

	var b strings.Builder

	// Header
	if m.filter != "" {
		fmt.Fprintf(&b, "\n  %spressctl%s  %s%s%s   %s/ %s_%s\n\n",
			bold, reset, dim, repoURL, reset, yellow, m.filter, reset)
	} else {
		fmt.Fprintf(&b, "\n  %spressctl%s  %s%s%s\n\n", bold, reset, dim, repoURL, reset)
	}

	// Items (only visible ones)
	if len(m.visible) == 0 {
		fmt.Fprintf(&b, "  %sno matches%s\n", dim, reset)
	} else {
		for pos, idx := range m.visible {
			n := pos + 1
			if pos == m.cursor {
				fmt.Fprintf(&b, "  %s▶ %-2d %s%s\n", cyan, n, m.items[idx], reset)
			} else {
				fmt.Fprintf(&b, "    %s%-2d %s%s\n", dim, n, m.items[idx], reset)
			}
		}
	}

	// Footer
	if m.filter != "" {
		fmt.Fprintf(&b, "\n  %s↑↓ navigate  enter select  esc clear filter%s\n", dim, reset)
	} else {
		fmt.Fprintf(&b, "\n  %s↑↓ navigate  enter or number to select  type to filter%s\n", dim, reset)
	}

	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "press",
	Short: "pressctl - Ansible wrapper for WordPress hosting management",
	Long: `pressctl is a CLI tool that simplifies WordPress hosting management
by wrapping Ansible playbooks with an intuitive, interactive interface.

Manage servers, sites, and domains with ease while maintaining full
visibility into your infrastructure state via ~/.pressctl/pressctl.yaml

Examples:
  # Add and provision a new server
  press server provision

  # Create a WordPress site
  press site create

  # Add a domain with SSL
  press domain add

  # List all servers
  press server list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		type action struct {
			label string
			run   func()
		}

		actions := []action{
			{"Provision a new server", func() { serverProvisionCmd.Run(serverProvisionCmd, []string{}) }},
			{"Create a WordPress site", func() { siteCreateCmd.Run(siteCreateCmd, []string{}) }},
			{"Delete a site", func() { siteDeleteCmd.Run(siteDeleteCmd, []string{}) }},
			{"Add a domain to a site", func() { domainAddCmd.Run(domainAddCmd, []string{}) }},
			{"Issue / renew SSL for a domain", func() { domainSSLCmd.Run(domainSSLCmd, []string{}) }},
			{"Check server health", func() { serverHealthCheckCmd.Run(serverHealthCheckCmd, []string{}) }},
			{"List servers", func() { serverListCmd.Run(serverListCmd, []string{}) }},
			{"List sites", func() { siteListCmd.Run(siteListCmd, []string{}) }},
		}

		labels := make([]string, len(actions))
		for i, a := range actions {
			labels[i] = a.label
		}

		result, err := tea.NewProgram(newMenuModel(labels)).Run()
		if err != nil {
			os.Exit(1)
		}

		chosen := result.(menuModel).chosen
		if chosen < 0 {
			os.Exit(1)
		}

		fmt.Println()
		actions[chosen].run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&DryRun, "dry-run", false, "Show what would be done without making changes")
}
