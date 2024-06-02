package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "waitdeps",
	Short: "Waits until dependent services are available",
	Long: `This CLI tool ensures that dependent services such as databases, SMTP
servers, or any other TCP-based services are available before proceeding with 
the application startup or tests. By waiting for these dependencies, it helps
prevent the application from starting with invalid configurations, which could
result in inconsistent states or failures.

This tool is particularly useful in production environments to guarantee that
all services are properly configured and ready, as well as in development
scenarios using tools like docker-compose. It ensures that all necessary
dependencies are up and running before the main application or tests are
executed. 

Example:
  waitdeps wait --tcp-connect postgres://user@db --tcp-connect 10.1.2.3:25`,
	Version: "0.0.1-alpha",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
