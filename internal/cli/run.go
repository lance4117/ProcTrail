package cli

import (
	"github.com/spf13/cobra"
)

// runCmd represents the hello command
var runCmd = &cobra.Command{
	Use:   "start [pid/name]",
	Short: "start to monitor process [pid/name] ",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			println("no target")
			return
		} else {
			println("target is ", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
