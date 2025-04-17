package cmd

import (
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "community-invite",
	Short: "Generate, test and send community meeting invitations",
	Long:  "CLI tool for managing OCM community meeting communications",
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file (default is ./config.yaml)")

	// Register subcommands
	rootCmd.AddCommand(NewGenerateCmd())
	rootCmd.AddCommand(NewTestmailCmd())
	rootCmd.AddCommand(NewSendCmd())
}
