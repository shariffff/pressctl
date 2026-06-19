package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/pressctl/cli/internal/web"
)

var (
	uiPort   int
	uiNoOpen bool
)

// uiCmd launches a local web dashboard rendering servers and sites.
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Open the web dashboard",
	Long: `Start a local web server that renders your pressctl configuration
(servers, sites, and domains) as a dashboard in your browser.

The dashboard is read-only and reflects the current contents of
~/.pressctl/pressctl.yaml. It reloads on every page refresh.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, _ := ensureConfig()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		if err := web.Serve(ctx, mgr, uiPort, !uiNoOpen); err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().IntVar(&uiPort, "port", 0, "Port to listen on (0 = random free port)")
	uiCmd.Flags().BoolVar(&uiNoOpen, "no-open", false, "Do not open the browser automatically")
}
