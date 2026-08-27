package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"certsmanager/src/internal/pki"
	"github.com/spf13/cobra"
)

// initCACmd handles the initialization of a new CA for a project.
var initCACmd = &cobra.Command{
	Use:   "init-ca [project_name]",
	Short: "Initializes the CA infrastructure for a specific project.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		fmt.Printf("Starting the CA initialization process for project: %s...\n", projectName)

		// Define the base directory for the project PKI
		projectDir := filepath.Join(".", "pki", projectName)

		// 1. Execute the CA initialization logic
		if err := pki.InitializeCA(projectName, projectDir); err != nil {
			fmt.Printf("ERROR initializing CA for %s: %v\n", projectName, err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Success! The CA has been initialized and the files (ca.key, ca.crt) are available at:")
		fmt.Printf("   %s/ca.crt\n", projectDir)
		fmt.Printf("   %s/ca.key (Keep it secret!)\n", projectDir)
		fmt.Println("Next step: Use 'certsmanager add-server <project> <server>' to issue server certificates.")
	},
}

func init() {
	rootCmd.AddCommand(initCACmd)
}

