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
// <keyPrefix>_<variable> = <placeholderPrefix>_<variable>
//
// actual environment
// <variable> = <value>

type Entry struct {
	key               string
	placeholder       string
	keyPrefix         string
	placeholderPrefix string
}

func NewEntry(key, placeholder, keyPrefix, placeholderPrefix string) Entry {
	return Entry{
		key:               key,
		placeholder:       placeholder,
		keyPrefix:         utils.AddSuffix(keyPrefix, "_"),
		placeholderPrefix: utils.AddSuffix(placeholderPrefix, "_"),
	}
}

func (e Entry) Key() string {
	return e.key
}

func (e Entry) Placeholder() string {
	return e.placeholder
}

func (e Entry) GeneratePlaceholder() string {
	return fmt.Sprintf("%s%s", e.placeholderPrefix, e.GetEnvVariable())
}

func (e Entry) GetEnvVariable() string {
	return strings.TrimPrefix(e.key, e.keyPrefix)
}

func (e Entry) GetEnvValue() string {
	return os.Getenv(e.GetEnvVariable())
}

func (e Entry) Skip() bool {
	return !strings.HasPrefix(e.key, e.keyPrefix)
}

func (e Entry) Validate() error {
	if strings.TrimPrefix(e.Key(), e.keyPrefix) != strings.TrimPrefix(e.Placeholder(), e.placeholderPrefix) {
		return ErrValidationFailed
	}
	return nil
}
