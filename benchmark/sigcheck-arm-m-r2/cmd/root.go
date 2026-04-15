package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sigcheck",
	Short: "sigcheck verifies file integrity against SHA-256 manifests",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
