package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"certsmanager/src/internal/pki"
	"github.com/spf13/cobra"
)

// caCmd represents the init-ca command
var caCmd = &cobra.Command{
	Use:   "init-ca [projectName]",
	Short: "Initializes the CA structure for a specific project.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		projectDir := filepath.Join("pki", projectName)

		fmt.Printf("\n--- Initializing CA for project: %s ---\n", projectName)
		
		if err := pki.InitializeCA(projectName, projectDir); err != nil {
			fmt.Printf("Error initializing CA: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n✅ Successfully initialized and generated CA key/certificate for '%s' in %s/\n", projectName, projectDir)
	},
}

func init() {
	rootCmd.AddCommand(caCmd)
}

