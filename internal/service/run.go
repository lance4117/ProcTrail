package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the hello command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long: `
run this app`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proctrail is running")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
