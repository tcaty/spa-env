package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapDotenvToActualEnv(t *testing.T) {
	testCases := []struct {
		name         string
		dotenvData   string
		actualEnv    map[string]string
		excpectedMap map[string]string
		fail         bool
	}{
		{
			name: "Valid dotenv file",
			dotenvData: `
VITE_API_URL=API_URL
VITE_SECRET_TOKEN=SECRET_TOKEN
			`,
			actualEnv: map[string]string{
				"API_URL":      "https://api.com/",
				"SECRET_TOKEN": "12345",
			},
			excpectedMap: map[string]string{
				"API_URL":      "https://api.com/",
				"SECRET_TOKEN": "12345",
			},
		},
		{
			name:       "Env variable specified in dotenv file but missed in env",
			dotenvData: "VITE_SECRET_TOKEN=SECRET_TOKEN",
			actualEnv: map[string]string{
				"API_URL": "https://api.com/",
			},
			fail: true,
		},
		{
			name: "Empty .env file",
			fail: true,
		},
		{
			name:       "Empty actual env",
			dotenvData: "VITE_API_URL=API_URL",
			fail:       true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			path := "/tmp/.env"

			prepareDotenv(t, path, []byte(tc.dotenvData))
			prepareEnv(t, tc.actualEnv)

			envMap, err := MapDotenvToActualEnv(path)
			if !tc.fail {
				require.NoError(t, err)
			}

			for k := range tc.excpectedMap {
				require.Equal(t, tc.excpectedMap[k], envMap[k])
			}
		})
	}
}

func prepareDotenv(t *testing.T, path string, data []byte) {
	t.Cleanup(func() {
		os.Remove(path)
	})

	f, err := os.Create(path)
	require.NoError(t, err)

	_, err = f.Write(data)
	require.NoError(t, err)
}

func prepareEnv(t *testing.T, env map[string]string) {
	t.Cleanup(func() {
		os.Clearenv()
	})

	for k, v := range env {
		err := os.Setenv(k, v)
		require.NoError(t, err)
	}
}
