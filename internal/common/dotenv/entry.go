package dotenv

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/tcaty/spa-env/internal/common/utils"
)

var ErrValidationFailed = errors.New(".env entry validation failed")

// .env files
// <      -- key --     >   <      -- placeholder --     >
// (<keyPrefix>_)<variable> = (<placeholderPrefix>_)<variable>
//
// actual environment
// <variable> = <value>

type Entry struct {
	Key               string
	Value             string
	KeyPrefix         string
	PlaceholderPrefix string
}

func NewEntry(key, value, keyPrefix, placeholderPrefix string) Entry {
	return Entry{
		Key:               key,
		Value:             value,
		KeyPrefix:         utils.AddSuffix(keyPrefix, "_"),
		PlaceholderPrefix: utils.AddSuffix(placeholderPrefix, "_"),
	}
}

func (e Entry) Placeholder() string {
	return fmt.Sprintf("%s%s", e.PlaceholderPrefix, e.EnvVariable())
}

func (e Entry) EnvVariable() string {
	return strings.TrimPrefix(e.Key, e.KeyPrefix)
}

func (e Entry) EnvValue() string {
	return os.Getenv(e.EnvVariable())
}

func (e Entry) Skip() bool {
	return !strings.HasPrefix(e.Key, e.KeyPrefix)
}

func (e Entry) Validate() error {
	if strings.TrimPrefix(e.Key, e.KeyPrefix) != strings.TrimPrefix(e.Placeholder(), e.PlaceholderPrefix) {
		return ErrValidationFailed
	}
	return nil
}
