package cmd

import (
	"log"

	"github.com/nickdala/azure-resource-verifier/internal/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the Azure Resource Verifier MCP server",
	Long:  `The mcp command starts the Azure Resource Verifier MCP server.`,
	Run:   mcpCommand,
}

func mcpCommand(cmd *cobra.Command, args []string) {
	log.Println("Starting MCP server...")
	logFile := viper.GetString("log-file")

	// Create a new resource verifier server
	resourceVerifierServer, err := mcp.NewResourceVerifierServer(logFile)
	if err != nil {
		log.Fatalf("Error creating resource verifier server: %v\n", err)
	}

	// Start the server with the default transport (stdio)
	if err := resourceVerifierServer.ServeStdio(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
