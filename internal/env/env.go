package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Map .env file to actual environment, return error
// if variable specified in file, but missed in current env
func MapEnvFileToActualEnv(path string) (map[string]string, error) {
	dotenv, err := godotenv.Read(path)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading .env file: %v", err)
	}

	res := make(map[string]string)

	for k, v := range dotenv {
		env := os.Getenv(v)

		if env == "" {
			return nil, fmt.Errorf("env variable %s specified in .env file, but not found in environment", k)
		} else {
			res[v] = env
		}

	}

	return res, nil
}
