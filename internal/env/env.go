package env

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/tcaty/spa-env/internal/log"
)

var (
	ErrMissedVariable = errors.New("env variable key specified in file, but placeholder missed in environment")
)

// Map .env file to actual environment, return error
// if variable specified in file, but missed in current environment
func MapDotenvToActualEnv(path string, prefix string) (map[string]string, error) {
	dotenv, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading %s file: %v", path, err)
	}

	res := make(map[string]string)

	for k, v := range dotenv {
		if !strings.HasPrefix(k, prefix) {
			log.Debug(
				"skip variable cause it hasn't prefix",
				"key", k,
				"prefix", prefix,
				"path", path,
			)
			continue
		}

		env := os.Getenv(v)

		if env == "" {
			log.Error(
				"missed variable",
				ErrMissedVariable,
				"key", k,
				"placeholder", v,
				"path", path,
			)
			err = ErrMissedVariable
		}

		res[v] = env
	}

	if err == ErrMissedVariable {
		return nil, fmt.Errorf("some variables from %s wasn't found in environment", path)
	}

	return res, nil
}
