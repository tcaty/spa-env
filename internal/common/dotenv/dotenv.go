package dotenv

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tcaty/spa-env/internal/common/log"
	"github.com/tcaty/spa-env/internal/common/utils"
	"github.com/tcaty/spa-env/pkg/file"
)

// Find dotenv file in workdir by filename and read it
// return map in form of [key]: [placeholder]
func Read(workdir, filename, keyPrefix, placeholderPrefix string) ([]Entry, error) {
	path, err := file.Find(workdir, filename)
	if err != nil {
		return nil, fmt.Errorf("error occured while finding .env file: %v", err)
	}

	content, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading %s file: %v", path, err)
	}

	entries := ParseEntries(content, keyPrefix, placeholderPrefix)

	log.Debug(
		".env file was read successfully",
		"path", path,
	)

	return entries, nil
}

func ParseEntries(dotenvContent map[string]string, keyPrefix, placeholderPrefix string) []Entry {
	// TODO: optimize
	entries := make([]Entry, 0)

	for key, value := range dotenvContent {
		entry := NewEntry(key, value, keyPrefix, placeholderPrefix)
		entries = append(entries, entry)
	}

	return entries
}

// TODO: move this one to pkg/file/file.go, replace entries and
// Find dotenv file in workdir by filename and write envMap there
func Write(workdir, filename string, entries []Entry, enableComments bool) error {
	path := fmt.Sprintf("%s%s", utils.AddSuffix(workdir, "/"), filename)

	content := generateContent(entries, enableComments)
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(content + "\n")
	if err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	log.Debug(
		"successful write to .env file",
		"path", path,
	)

	return nil
}

// TODO: move this one to internal/generate/generate.go as a core logic
func generateContent(entries []Entry, enableComments bool) string {
	content := make([]string, 0)

	if enableComments {
		content = append(
			content,
			"# This file was auto-generated by spa-env tool. Don't edit it manually!",
			"# There is a full list of client environment variables below",
			"#",
		)

		for _, entry := range entries {
			content = append(content, fmt.Sprintf("# %s", entry.EnvVariable()))
		}
	}

	for _, entry := range entries {
		if entry.Skip() {
			continue
		}

		if enableComments {
			content = append(
				content,
				fmt.Sprintf("\n# env -> %s", entry.EnvVariable()),
				fmt.Sprintf("# src -> process.env.%s", entry.Key),
			)
		}

		content = append(content, fmt.Sprintf("%s=%s", entry.Key, entry.Placeholder()))
	}

	return strings.Join(content, "\n")
}
