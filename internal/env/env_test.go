package env

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapDotenvToActualEnv(t *testing.T) {
	testCases := []struct {
		name         string
		dotenvData   string
		prefix       string
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
		{
			name: "Use env prefix",
			dotenvData: `
POSTGRES_CONN_STRING=POSTGRES_CONN_STRING
NEXT_PUBLIC_API_URL=API_URL
NEXT_PUBLIC_TOKEN=SECRET_TOKEN
			`,
			prefix: "NEXT_PUBLIC",
			actualEnv: map[string]string{
				"POSTGRES_CONN_STRING": "postgres://username:password@localhost:5432/database",
				"API_URL":              "https://api.com/",
				"SECRET_TOKEN":         "12345",
			},
			excpectedMap: map[string]string{
				"API_URL":      "https://api.com/",
				"SECRET_TOKEN": "12345",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			path := "/tmp/.env"

			prepareDotenv(t, path, []byte(tc.dotenvData))
			prepareEnv(t, tc.actualEnv)

			envMap, err := MapDotenvToActualEnv(path, tc.prefix, false)
			if !tc.fail {
				require.NoError(t, err)
			}

			if !tc.fail && !reflect.DeepEqual(tc.excpectedMap, envMap) {
				require.Failf(t, "Maps should be equal", "Maps aren't equal \nExpected: %#v \nActual: %#v", tc.excpectedMap, envMap)
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
