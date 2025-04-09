package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/y7ut/mcp-tavily-search/internal/tavily"
	"github.com/y7ut/mcp-tavily-search/internal/tool"
)

// debug flag
var debug bool

// RunCmd
// environment variables:
// TRVILY_API_KEY = "your tavily api key"
// TRVILY_INCLUDE_DOMAINS = "domain1,domain2"
// TRVILY_EXCLUDE_DOMAINS = "domain1,domain2"
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		trvilyApiKey, _ := os.LookupEnv("TRVILY_API_KEY")
		if len(args) > 0 {
			trvilyApiKey = args[0]
		}
		if trvilyApiKey == "" {
			fmt.Println("TRVILY_API_KEY is required")
			os.Exit(1)
		}
		var includeDomain []string
		if os.Getenv("TRVILY_INCLUDE_DOMAINS") != "" {
			includeDomain = strings.Split(os.Getenv("TRVILY_INCLUDE_DOMAINS"), ",")
		}

		var excludeDomain []string
		if os.Getenv("TRVILY_EXCLUDE_DOMAINS") != "" {
			excludeDomain = strings.Split(os.Getenv("TRVILY_EXCLUDE_DOMAINS"), ",")
		}

		tavily.Init(trvilyApiKey, debug, includeDomain, excludeDomain)
		mcpServerRun()
	},
}

func init() {
	RootCmd.AddCommand(RunCmd)

	RunCmd.Flags().BoolVarP(&debug, "debug", "d", true, "Enable debug mode")
}

// mcpServerRun run the mcp server
func mcpServerRun() {
	// Create MCP server
	s := server.NewMCPServer(
		"MCP Tavily Search 🔍",
		"1.0.0",
		server.WithLogging(),
		server.WithResourceCapabilities(true, true),
	)

	tool.Bind(s)
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
