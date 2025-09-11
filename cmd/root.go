package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	activityMonitor "github.com/tribock/activitymonitor/pkg/activity_monitor"
)

var (
	stats   bool
	timeout int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "activitymonitor",
	Short: "ActivityMonitor keeps your computer awake by monitoring and simulating activity",
	Long: `ActivityMonitor is a Go program focused on monitoring user activity and collecting 
statistics about computer usage. It helps keep your screen active to avoid session 
timeouts or screen lock by periodically moving the mouse cursor by 1 pixel and then 
moving it back. If manual mouse movement is detected, the program records this as user 
activity and skips the automated movement for that cycle.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create the activity monitor
		monitor := activityMonitor.NewActivityMonitor()

		// Add stats if flag is set
		if stats {
			monitor = monitor.WithStats()
		}

		// Add timeout if flag is set
		if timeout > 0 {
			monitor = monitor.WithTimeout(time.Duration(timeout) * time.Second)
		}

		// If no flags are set, use default behavior (stats enabled, 1 minute timeout)
		if !stats && timeout == 0 {
			monitor = monitor.WithStats().WithTimeout(1 * time.Minute)
		}

		// Start monitoring
		monitor.KeepOnMoving()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Define flags
	rootCmd.Flags().BoolVarP(&stats, "stats", "s", false, "Enable statistics collection and reporting")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 0, "Set timeout interval in seconds (default: 60 seconds when no flags are provided)")

	// Add help text for flags
	rootCmd.Flags().Lookup("stats").Usage = "Enable statistics collection and reporting about user activity"
	rootCmd.Flags().Lookup("timeout").Usage = "Set the interval between activity checks in seconds (e.g., -t 30 for 30 seconds)"
}
