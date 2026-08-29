package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"certsmanager/src/internal/pki"
	"github.com/spf13/cobra"
)

// addServerCmd handles the issuance of server certificates.
var addServerCmd = &cobra.Command{
	Use:   "add-server [project_name] [server_name]",
	Short: "Generates a server certificate using a project's CA.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		serverName := args[1]
		
		// Corrected: Handling the error from GetString
		sansList, err := cmd.Flags().GetString("sans")
		if err != nil {
			fmt.Printf("Error reading 'sans' flag: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Starting certificate issuance for project '%s' (Server: %s)...\n", projectName, serverName)

		// 1. Determine project directory
		projectDir := filepath.Join(".", "pki", projectName)

		// Check if CA exists
		if _, err := os.Stat(filepath.Join(projectDir, "ca.crt")); os.IsNotExist(err) {
			fmt.Printf("ERROR: CA not found for project '%s'. Please run 'certsmanager init-ca %s' first.\n", projectName, projectName)
			os.Exit(1)
		}

		// 2. Execute certificate generation logic
		if err := pki.AddServerCertificate(projectName, serverName, sansList, projectDir); err != nil {
			fmt.Printf("ERROR generating certificate for %s: %v\n", serverName, err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Success! The server certificate has been generated and saved at:")
		fmt.Printf("   %s/servers/%s_server.crt\n", projectDir, serverName)
		fmt.Println("Remember to manage private keys securely.")
	},
}

func init() {
	// Flag for SANs (Subject Alternative Names)
	addServerCmd.Flags().StringP("sans", "s", "", "Comma-separated list of SANs (e.g., dns1,dns2,ip1)")
	// Add to rootCmd
	rootCmd.AddCommand(addServerCmd)
}

