package dotenv

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/tcaty/spa-env/internal/common/utils"
)

var ErrValidationFailed = errors.New(".env entry validation failed")

// Entry is a CORE struct in the current project. It is used for working with
// entries in dotenv files as well as for working with environment variables.
//
// Despite the conditions exactly this struct must be used for the purposes above
// in order to keep the same terms and principles through entire project.
//
// ** CORE TERMS **
//
// Place: .env files
// <        -- Key --        >   <     -- Value/Placeholder --     >
// (<keyPrefix>_)<EnvVariable> = (<placeholderPrefix>_)<EnvVariable>
// ! Understand distinction between value and placehodler:
// value - not validated value, that was just parsed from .env file
// it might has some random value or already has a placeholder.
// placeholder - validated value that corresponds particular format above.
//
// Place: actual environment
// <EnvVariable> = <EnvValue>
//
// **
type Entry struct {
	// left side "=" in .env file
	Key string
	// right side "=" in .env file
	Value string
	// key prefix from cmd flags
	KeyPrefix string
	// placeholder prefix from cmd flags
	PlaceholderPrefix string
}

func NewEntry(key, value, keyPrefix, placeholderPrefix string) Entry {
	return Entry{
		Key:   key,
		Value: value,
		// add "_" suffix in order to trim these prefixes correctly
		KeyPrefix:         utils.AddSuffix(keyPrefix, "_"),
		PlaceholderPrefix: utils.AddSuffix(placeholderPrefix, "_"),
	}
}

// Environment Variable is set in the actual environment see  ** CORE TERMS **.
// This variables is used docker-compose.yml files environment.
func (e Entry) EnvVariable() string {
	return strings.TrimPrefix(e.Key, e.KeyPrefix)
}

// Environment Value is set in the actual environment see  ** CORE TERMS **.
// This value is taken exactly from Environment Variable and
// it is used in order to replace placeholder in built apps.
func (e Entry) EnvValue() string {
	return os.Getenv(e.EnvVariable())
}

// Placeholder is value that replaces environment variables in static files.
// See root README.md to understand how does env varialbe work in spa.
func (e Entry) Placeholder() string {
	return fmt.Sprintf("%s%s", e.PlaceholderPrefix, e.EnvVariable())
}

// Check that entry key has necessarry prefix, skip if it's not.
// For example, in NextJS apps server side environments variables mustn't be affected.
// Also this kind of variables has no any prefix, therefore we must skip it.
func (e Entry) Skip() bool {
	return !strings.HasPrefix(e.Key, e.KeyPrefix)
}

// Check that entry corresponds to the scheme from ** CORE TERMS **.
func (e Entry) Validate() error {
	if strings.TrimPrefix(e.Key, e.KeyPrefix) != strings.TrimPrefix(e.Placeholder(), e.PlaceholderPrefix) {
		return ErrValidationFailed
	}
	return nil
}
