package replace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tcaty/spa-entrypoint/internal/env"
	"github.com/tcaty/spa-entrypoint/internal/file"
)

// Form replacement rules by parsing dotenv file and actual env
// Then recursively walk through files in workdir
// and replace env variables by formed rules
func Replace(workdir string, dotenv string, verbose bool) error {
	if _, err := os.Stat(workdir); err != nil {
		return fmt.Errorf("error occured while reading workdir: %v", err)
	}

	dotenvPath, err := file.Find(workdir, dotenv)
	if err != nil {
		return fmt.Errorf("error occured while finding .env file: %v", err)
	}

	// form replacement rules
	rules, err := env.MapDotenvToActualEnv(dotenvPath)
	if err != nil {
		return fmt.Errorf("error occured while mapping .env file to env: %v", err)
	}

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

		if err := file.Replace(path, rules, verbose); err != nil {
			return fmt.Errorf("error occured while replacing file content: %v", err)
		}

		return nil
	})

	return err
}
