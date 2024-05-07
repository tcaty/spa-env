package replace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tcaty/spa-env/internal/env"
	"github.com/tcaty/spa-env/internal/file"
)

// Form replacement rules by parsing dotenv file and actual env
// Then recursively walk through files in workdir
// and replace env variables by formed rules
// Return updated files count and no error if replacement completed successfully
func Replace(workdir string, dotenv string, prefix string) (int, error) {
	if _, err := os.Stat(workdir); err != nil {
		return 0, fmt.Errorf("error occured while reading workdir: %v", err)
	}

	dotenvPath, err := file.Find(workdir, dotenv)
	if err != nil {
		return 0, fmt.Errorf("error occured while finding .env file: %v", err)
	}

	// form replacement rules
	rules, err := env.MapDotenvToActualEnv(dotenvPath, prefix)
	if err != nil {
		return 0, fmt.Errorf("error occured while mapping .env file to env: %v", err)
	}

	filesUpdated := 0

	err = filepath.WalkDir(workdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("prevent panic by handling failure accessing a path %q: %v", path, err)
		}

		// skip entire node_modules directory
		if d.IsDir() && d.Name() == "node_modules" {
			return filepath.SkipDir
		}

		// skip directories paths and dotenv file
		if d.IsDir() || strings.Contains(path, dotenv) {
			return nil
		}

		updated, err := file.ReplaceContent(path, rules)
		if err != nil {
			return fmt.Errorf("error occured while replacing file content: %v", err)
		}

		if updated {
			filesUpdated += 1
		}

		return nil
	})

	return filesUpdated, err
}
