package env

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// TODO: add docs with examples
func MapEnvFileToActualEnv(path string) (map[string]string, error) {
	env, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading .env file: %v", err)
	}

	res := make(map[string]string)

	for k, v := range env {
		envVar := os.Getenv(v)

		if envVar == "" {
			slog.Warn(
				"variable specified in file, but not specefied in current env",
				"var", k,
			)
		} else {
			res[v] = envVar
		}

	}

	return res, nil
}
