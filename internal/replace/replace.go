package replace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tcaty/spa-env/internal/log"
	"github.com/tcaty/spa-env/pkg/file"
)

var (
	ErrMissedVariable = errors.New("env variable key specified in file, but placeholder missed in environment")
)

// Form replacement rules by parsing dotenv file and actual env
// Then recursively walk through files in workdir
// and replace env variables by formed rules
// Return updated files count and no error if replacement completed successfully
func Replace(workdir string, dotenv string, keyPrefix string, placeholderPrefix string) (int, error) {
	if _, err := os.Stat(workdir); err != nil {
		return 0, fmt.Errorf("error occured while reading workdir: %v", err)
	}

	dotenvPath, err := file.Find(workdir, dotenv)
	if err != nil {
		return 0, fmt.Errorf("error occured while finding .env file: %v", err)
	}

	// form replacement rules
	rules, err := mapPlaceholderToValue(dotenvPath, keyPrefix, placeholderPrefix)
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

		// skip directories paths and all dotenv files
		if d.IsDir() || strings.HasPrefix(d.Name(), ".env") {
			return nil
		}

		appliedRules, err := file.ReplaceContent(path, rules)
		if err != nil {
			return fmt.Errorf("error occured while replacing file content: %v", err)
		}

		if len(appliedRules) > 0 {
			filesUpdated += 1
			for placeholder, value := range appliedRules {
				log.Debug(
					"Successful replacement",
					"path", path,
					"placeholder", placeholder,
					"value", value,
				)
			}
		}

		return nil
	})

	return filesUpdated, err
}

// Map placeholders from .env file to actual environment variables values
// return error if variable specified in file, but missed in current environment
func mapPlaceholderToValue(path string, keyPrefix string, placeholderPrefix string) (map[string]string, error) {
	dotenv, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading %s file: %v", path, err)
	}

	res := make(map[string]string)

	for key, placeholder := range dotenv {
		if !strings.HasPrefix(key, keyPrefix) {
			log.Debug(
				"Skip variable cause it has no prefix",
				"key", key,
				"prefix", keyPrefix,
				"path", path,
			)
			continue
		}

		value := getenv(placeholder, placeholderPrefix)

		if value == "" {
			log.Error(
				"missed variable",
				ErrMissedVariable,
				"key", key,
				"placeholder", placeholder,
				"path", path,
			)
			err = ErrMissedVariable
		}

		res[placeholder] = value
	}

	if err == ErrMissedVariable {
		return nil, fmt.Errorf("some variables from %s wasn't found in environment", path)
	}

	return res, nil
}

// Find environment variable without prefix by key with prefix
// For example:
// key := PLACEHOLDER_TOKEN, prefix := PLACEHOLDER -> value from env variable TOKEN
// key := PLACEHOLDER_TOKEN, prefix := "" -> value from env variable PLACEHOLDER_TOKEN
func getenv(key string, prefix string) string {
	if prefix == "" {
		return os.Getenv(key)
	}

	// add "_" to the end of the prefix if there is no it yet
	if !strings.HasSuffix(prefix, "_") {
		prefix = fmt.Sprintf("%s_", prefix)
	}

	key = strings.TrimPrefix(key, prefix)

	return os.Getenv(key)
}
