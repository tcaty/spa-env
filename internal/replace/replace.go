package replace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tcaty/spa-env/internal/common/dotenv"
	"github.com/tcaty/spa-env/internal/common/log"
	"github.com/tcaty/spa-env/pkg/file"
)

var (
	errMissedVariable = errors.New("env variable key specified in file, but placeholder missed in environment")
)

// Form replacement rules by parsing dotenv file and actual env
// Then recursively walk through files in workdir
// and replace env variables by formed rules
// Return updated files count and no error if replacement completed successfully
func Replace(workdir, dotenvProd, keyPrefix, placeholderPrefix string) (int, error) {
	if _, err := os.Stat(workdir); err != nil {
		return 0, fmt.Errorf("error occured while reading workdir: %v", err)
	}

	dotenvContent, err := dotenv.Read(workdir, dotenvProd)
	if err != nil {
		return 0, fmt.Errorf("error occured while reading .env file: %v", err)
	}

	// form replacement rules
	rules, err := mapPlaceholderToValue(dotenvContent, keyPrefix, placeholderPrefix)
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
func mapPlaceholderToValue(dotenv map[string]string, keyPrefix string, placeholderPrefix string) (map[string]string, error) {
	var err error
	res := make(map[string]string)

	for key, placeholder := range dotenv {
		if !strings.HasPrefix(key, keyPrefix) {
			log.Debug(
				"skip variable cause it has no prefix",
				"key", key,
				"prefix", keyPrefix,
			)
			continue
		}

		value := getenv(placeholder, placeholderPrefix)

		if value == "" {
			log.Error(
				"missed variable",
				errMissedVariable,
				"key", key,
				"placeholder", placeholder,
			)
			err = errMissedVariable
		}

		res[placeholder] = value
	}

	if errors.Is(err, errMissedVariable) {
		return nil, err
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

	key = strings.TrimPrefix(key, prefix)

	return os.Getenv(key)
}
