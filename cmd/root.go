package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "mcp-tavily-search",
	Short:   "A server implement Model Context Protocol, used for searching from tavily.",
	Version: "1.0.0",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
