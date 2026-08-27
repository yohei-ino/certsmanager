package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"certsmanager/src/internal/pki"
	"github.com/spf13/cobra"
)

// serverCmd represents the add-server command
var serverCmd = &cobra.Command{
	Use:   "add-server [projectName] [serverName]",
	Short: "Generates and signs a new server certificate for a project.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		serverName := args[1]
		
		// 修正: GetStringは (string, error) を返すためエラーハンドリングを追加
		sans, err := cmd.Flags().GetString("sans")
		if err != nil {
			fmt.Printf("Error reading 'sans' flag: %v\n", err)
			os.Exit(1)
		}

		projectDir := filepath.Join("pki", projectName)

		fmt.Printf("\n--- Generating server certificate for '%s' in project '%s' ---\n", serverName, projectName)
		
		// 修正: 関数名を AddServerCertificate に合わせ、引数を4つ渡す
		if err := pki.AddServerCertificate(projectName, serverName, sans, projectDir); err != nil {
			fmt.Printf("Error generating server certificate: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n✅ Successfully generated and saved server certificate for '%s' in %s/servers/\n", serverName, projectDir)
	},
}

func init() {
	// Define the --sans flag for the server command
	serverCmd.Flags().StringP("sans", "s", "", "Comma-separated list of Subject Alternative Names (e.g. example.com,www.example.com)")
	// Add the server command to the root command
	rootCmd.AddCommand(serverCmd)
}

