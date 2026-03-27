package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/pressctl/cli/internal/updater"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update press to the latest version",
	Long: `Check for a newer release on GitHub and, if found, replace the press
binary and Ansible playbooks in place. Your config (~/.pressctl/pressctl.yaml)
is never modified.

Examples:
  # Check and install if a newer version is available
  press update

  # Only check — print the latest version without installing
  press update --check`,
	Run: func(cmd *cobra.Command, args []string) {
		checkOnly, _ := cmd.Flags().GetBool("check")

		if Version == "dev" {
			color.Yellow("Running a development build — update skipped.")
			return
		}

		s := spinner.New(spinner.CharSets[14], 80*time.Millisecond)
		s.Suffix = "  Checking for updates..."
		s.Start()
		release, err := updater.LatestRelease()
		s.Stop()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		currentClean := strings.TrimPrefix(Version, "v")

		if !updater.IsNewer(Version, release.TagName) {
			color.Green("✓ Already on the latest version (%s)", currentClean)
			return
		}

		fmt.Printf("  Installed : v%s\n", currentClean)
		fmt.Printf("  Available : %s\n\n", release.TagName)

		if checkOnly {
			fmt.Printf("Run 'press update' to install.\n")
			return
		}

		s = spinner.New(spinner.CharSets[14], 80*time.Millisecond)
		s.Suffix = fmt.Sprintf("  Downloading %s...", release.TagName)
		s.Start()
		err = updater.Install(release.TagName)
		s.Stop()

		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Updated to %s", release.TagName)
		fmt.Println("  Run 'press version' to confirm.")
	},
}

func init() {
	updateCmd.Flags().Bool("check", false, "Check for a newer version without installing")
	rootCmd.AddCommand(updateCmd)
}
