package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/had-nu/sigcheck/internal/hasher"
	"github.com/had-nu/sigcheck/internal/manifest"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify files against a SHA-256 manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestPath, err := cmd.Flags().GetString("manifest")
		if err != nil {
			return err
		}
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			return err
		}

		if manifestPath == "" || target == "" {
			return fmt.Errorf("both --manifest and --target flags are required")
		}

		entries, err := manifest.Parse(manifestPath, target)
		if err != nil {
			return err
		}

		hasFailure := false
		for _, entry := range entries {
			fullPath := filepath.Join(target, entry.Path)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				fmt.Printf("MISS  %s\n", entry.Path)
				hasFailure = true
				continue
			}

			computed, err := hasher.HashFile(fullPath)
			if err != nil {
				return err
			}

			if computed != entry.Hash {
				fmt.Printf("FAIL  %s (hash mismatch)\n", entry.Path)
				hasFailure = true
			} else {
				fmt.Printf("OK    %s\n", entry.Path)
			}
		}

		if hasFailure {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().String("manifest", "", "Manifest file to verify against")
	verifyCmd.Flags().String("target", "", "Base directory for file resolution")
}
