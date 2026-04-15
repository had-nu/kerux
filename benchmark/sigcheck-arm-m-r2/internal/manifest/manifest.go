package manifest

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Entry struct {
	Hash string
	Path string
}

func Parse(manifestPath, baseDir string) ([]Entry, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var entries []Entry

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		parts := strings.Split(line, "  ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed line %d", lineNum)
		}

		hash, path := parts[0], parts[1]
		if len(hash) != 64 {
			return nil, fmt.Errorf("invalid hash length on line %d", lineNum)
		}

		if !safePath(path, baseDir) {
			return nil, fmt.Errorf("path traversal attempt: %s", path)
		}

		entries = append(entries, Entry{Hash: hash, Path: path})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func Write(outputPath string, entries []Entry) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, entry := range entries {
		if _, err := fmt.Fprintf(f, "%s  %s\n", entry.Hash, entry.Path); err != nil {
			return err
		}
	}
	return nil
}

func safePath(path, baseDir string) bool {
	// Reject paths containing directory traversal or absolute constructs early
	if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "..") || strings.Contains(path, "/../") {
		return false
	}
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return false
	}
	absTarget, err := filepath.Abs(filepath.Join(absBase, path))
	if err != nil {
		return false
	}
	return strings.HasPrefix(absTarget, absBase)
}
