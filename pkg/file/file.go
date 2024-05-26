package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Find file by workdir and filename substring
// Return relative path to file and no error
// Or return empty string and ErrNotExist if file wasn't found
func Find(wordkir string, filename string) (string, error) {
	var path string

	err := filepath.WalkDir(wordkir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("prevent panic by handling failure accessing a path %q: %v", p, err)
		}

		if !d.IsDir() && d.Name() == filename {
			path = p
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error occured while walking dir: %v", err)
	}

	if path == "" {
		return "", os.ErrNotExist
	}

	return path, nil
}

// Replace file content by provided rules
// Return appliedRules and no error if file was successfully updated
func ReplaceContent(path string, rules map[string]string) (map[string]string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading file: %v", err)
	}

	appliedRules := make(map[string]string)
	oldContent := string(bytes)
	newContent := oldContent

	for old, new := range rules {
		// skip if replacement isn't needed
		if !strings.Contains(newContent, old) {
			continue
		}

		newContent = strings.ReplaceAll(newContent, old, new)
		// mark current rule as applied
		appliedRules[old] = new
	}

	// prevent writing if there were no changes in content
	if oldContent == newContent {
		return nil, nil
	}

	if err := os.WriteFile(path, []byte(newContent), 0); err != nil {
		return nil, fmt.Errorf("error occured while writing file: %v", err)
	}

	return appliedRules, nil
}
