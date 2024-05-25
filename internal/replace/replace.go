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

// Parse dotenv file, validate parsed entries and form replacement rules
// Then recursively walk through files in workdir and replace env variables by formed rules
// Return updated files count and no error if replacement completed successfully
func Replace(workdir, dotenvProd, keyPrefix, placeholderPrefix string) (int, error) {
	if _, err := os.Stat(workdir); err != nil {
		return 0, fmt.Errorf("error occured while reading workdir: %v", err)
	}

	dotenvEntries, err := dotenv.Read(workdir, dotenvProd, keyPrefix, placeholderPrefix)
	if err != nil {
		return 0, fmt.Errorf("error occured while reading .env file: %v", err)
	}

	if err := validateEntries(dotenvEntries); err != nil {
		return 0, err
	}

	// form replacement rules
	rules, err := mapPlaceholderToValue(dotenvEntries)
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
					"successful replacement",
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

// Validate .env entries
func validateEntries(entries []dotenv.Entry) error {
	var validationErr error
	for _, entry := range entries {
		if entry.Skip() {
			continue
		}

		if err := entry.Validate(); err != nil {
			validationErr = err
			log.Error(
				".env entry validation failed", err,
				"excpected", fmt.Sprintf("%s=%s", entry.Key(), entry.GeneratePlaceholder()),
				"got", fmt.Sprintf("%s=%s", entry.Key(), entry.Placeholder()),
			)
		}
	}
	return validationErr
}

// Map placeholders from .env file to actual environment variables values
// return error if variable specified in file, but missed in current environment
func mapPlaceholderToValue(entries []dotenv.Entry) (map[string]string, error) {
	var err error
	res := make(map[string]string)

	for _, entry := range entries {
		if entry.Skip() {
			log.Debug(
				"skip variable cause it has no required prefix",
				"key", entry.Key(),
			)
			continue
		}

		value := entry.GetEnvValue()

		if value == "" {
			log.Error(
				"missed variable",
				errMissedVariable,
				"key", entry.Key(),
				"variable", entry.GetEnvVariable(),
			)
			err = errMissedVariable
		}

		res[entry.Placeholder()] = value
	}

	if errors.Is(err, errMissedVariable) {
		return nil, err
	}

	return res, nil
}
