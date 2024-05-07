package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tcaty/spa-env/internal/log"
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

// Replace placeholder by value in file located on specified path
func ReplaceContent(path string, rules map[string]string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error occured while reading file: %v", err)
	}

	oldContent := string(bytes)
	newContent := oldContent

	for placeholder, value := range rules {
		// skip if replacement isn't needed
		if !strings.Contains(newContent, placeholder) {
			continue
		}

		newContent = strings.ReplaceAll(newContent, placeholder, value)

		log.Debug(
			"successful replacement",
			"path", path,
			"placeholder", placeholder,
			"value", value,
		)
	}

	// prevent writing if there were no changes in content
	if oldContent == newContent {
		return nil
	}

	if err := os.WriteFile(path, []byte(newContent), 0); err != nil {
		return fmt.Errorf("error occured while writing file: %v", err)
	}

	return nil
}
