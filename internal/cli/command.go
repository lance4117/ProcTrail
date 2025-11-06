package cli

import (
	"os"

	"github.com/lance4117/gofuse/config"
	"github.com/spf13/cobra"
)

var (
	configPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "proctrail",
	Short: "\n ProcTrail — Record your process performance trail (CPU, memory, disk, network) to CSV, JSON.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		println(configPath)
		return config.Init(configPath)
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand()

	rootCmd.Flags().StringVar(&configPath, "config", "", "path to your config file. e.g. \"./path/config.yaml\" ")
}
