package dotenv

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/tcaty/spa-env/internal/common/log"
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
		".env file parsed successfully",
		"path", path,
	)

	return entries, nil
}

// TODO: sort in alphabetic order
func ParseEntries(dotenvContent map[string]string, keyPrefix, placeholderPrefix string) []Entry {
	entries := make([]Entry, 0, len(dotenvContent))

	for key, value := range dotenvContent {
		entry := NewEntry(key, value, keyPrefix, placeholderPrefix)
		entries = append(entries, entry)
	}

	return entries
}
