package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	// Flags
	pid          int32
	processName  string
	interval     int
	duration     int
	outputFormat string
	outputFile   string
	prettyPrint  bool
)

var rootCmd = &cobra.Command{
	Use:   "proctrail",
	Short: "\n ProcTrail - Process Metrics Collection Tool",
	Long: `
ProcTrail is a cross-platform CLI tool for collecting process metrics.
It can collect CPU usage, memory usage, I/O stats, thread count, and file descriptor count.`,
	RunE: runCommand,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Int32VarP(&pid, "pid", "p", 0, "Process ID to monitor")
	rootCmd.Flags().StringVarP(&processName, "name", "n", "", "Process name to monitor")
	rootCmd.Flags().IntVarP(&interval, "interval", "i", 0, "Sampling interval in seconds (0 for one-time sampling)")
	rootCmd.Flags().IntVarP(&duration, "duration", "d", 0, "Total duration to monitor in seconds (0 for indefinite)")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format (csv, json, jsonl, xml)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.Flags().BoolVar(&prettyPrint, "pretty", false, "Pretty print output (for json and xml)")
}

func validateFlags() error {
	if pid == 0 && processName == "" {
		return fmt.Errorf("either --pid or --name must be specified")
	}

	if pid != 0 && processName != "" {
		return fmt.Errorf("cannot specify both --pid and --name")
	}

	if interval < 0 {
		return fmt.Errorf("interval must be >= 0")
	}

	if duration < 0 {
		return fmt.Errorf("duration must be >= 0")
	}

	format := strings.ToLower(outputFormat)
	switch format {
	case "csv", "json", "jsonl", "xml":
		// valid formats
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	return nil
}

func getDuration() time.Duration {
	if duration > 0 {
		return time.Duration(duration) * time.Second
	}
	return 0
}
