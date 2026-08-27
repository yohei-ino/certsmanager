package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command for the 'certsmanager' CLI
var rootCmd = &cobra.Command{
	Use:   "certsmanager",
	Short: "A tool for managing PKI certificates in Go.",
	Long:  "CertsManager is a robust CLI tool to initialize, manage, and issue digital certificates for multiple projects, using Go's standard cryptographic library ecosystem.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the CertsManager CLI. Execute a subcommand to get started.")
	},
}

// Execute starts the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

