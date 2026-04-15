package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/had-nu/sigcheck/internal/hasher"
	"github.com/had-nu/sigcheck/internal/manifest"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a SHA-256 manifest from a target directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			return err
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		if target == "" || output == "" {
			return fmt.Errorf("both --target and --output flags are required")
		}

		var entries []manifest.Entry

		err = filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if d.Type()&os.ModeSymlink != 0 {
				return nil // T3: skip symlinks
			}

			relPath, err := filepath.Rel(target, path)
			if err != nil {
				return err
			}

			hash, err := hasher.HashFile(path)
			if err != nil {
				return err
			}

			entries = append(entries, manifest.Entry{Hash: hash, Path: relPath})
			return nil
		})

		if err != nil {
			return err
		}

		if err := manifest.Write(output, entries); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().String("target", "", "Directory to hash")
	generateCmd.Flags().String("output", "", "Manifest output path")
}
