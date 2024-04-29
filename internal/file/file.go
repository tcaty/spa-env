package file

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Find file by workdir and filename substring
// Return relative path to file and no error
// Or return empty string and ErrNotExist if file wasn't found
func Find(wordkir string, filenameSubstr string) (string, error) {
	var path string

	err := filepath.WalkDir(wordkir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("prevent panic by handling failure accessing a path %q: %v", p, err)
		}

		if strings.Contains(p, filenameSubstr) {
			path = p
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

// Replace old substring to new string in file located on specified path
func Replace(path string, rules map[string]string, verbose bool) error {
	read, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error occured while reading file: %v", err)
	}

	newContent := string(read)
	for old, new := range rules {
		newContent = strings.ReplaceAll(newContent, old, new)
		if verbose {
			log.Printf(
				"successfull replace path=%s, old=%s, new=%s",
				path, old, new,
			)
		}
	}

	if err := os.WriteFile(path, []byte(newContent), 0); err != nil {
		return fmt.Errorf("error occured while writing file: %v", err)
	}

	return nil
}
