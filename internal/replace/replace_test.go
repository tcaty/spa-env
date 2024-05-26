package replace

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tcaty/spa-env/internal/common/dotenv"
	"github.com/tcaty/spa-env/internal/common/log"
)

func TestMapPlaceholderToValue(t *testing.T) {
	testCases := []struct {
		name              string
		keyPrefix         string
		placeholderPrefix string
		dotenvContent     map[string]string
		actualEnv         map[string]string
		excpectedMap      map[string]string
		err               error
	}{
		{
			name: "Valid case without prefixes",
			dotenvContent: map[string]string{
				"API_URL":      "API_URL",
				"SECRET_TOKEN": "SECRET_TOKEN",
			},
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
			name:      "Valid case with keyPrefix",
			keyPrefix: "NEXT_PUBLIC_",
			dotenvContent: map[string]string{
				"POSTGRES_CONN_STRING":     "POSTGRES_CONN_STRING",
				"NEXT_PUBLIC_API_URL":      "API_URL",
				"NEXT_PUBLIC_SECRET_TOKEN": "SECRET_TOKEN",
			},
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
		{
			name:              "Valid case with placeholderPrefix",
			placeholderPrefix: "PLACEHOLDER_",
			dotenvContent: map[string]string{
				"API_URL":      "PLACEHOLDER_API_URL",
				"SECRET_TOKEN": "PLACEHOLDER_SECRET_TOKEN",
			},
			actualEnv: map[string]string{
				"API_URL":      "https://api.com/",
				"SECRET_TOKEN": "12345",
			},
			excpectedMap: map[string]string{
				"PLACEHOLDER_API_URL":      "https://api.com/",
				"PLACEHOLDER_SECRET_TOKEN": "12345",
			},
		},
		{
			name:              "Valid case with keyPrefix and placeholderPrefix",
			keyPrefix:         "NEXT_PUBLIC_",
			placeholderPrefix: "PLACEHOLDER_",
			dotenvContent: map[string]string{
				"POSTGRES_CONN_STRING":     "PLACEHOLDER_POSTGRES_CONN_STRING",
				"NEXT_PUBLIC_API_URL":      "PLACEHOLDER_API_URL",
				"NEXT_PUBLIC_SECRET_TOKEN": "PLACEHOLDER_SECRET_TOKEN",
			},
			actualEnv: map[string]string{
				"POSTGRES_CONN_STRING": "postgres://username:password@localhost:5432/database",
				"API_URL":              "https://api.com/",
				"SECRET_TOKEN":         "12345",
			},
			excpectedMap: map[string]string{
				"PLACEHOLDER_API_URL":      "https://api.com/",
				"PLACEHOLDER_SECRET_TOKEN": "12345",
			},
		},
		{
			name:          "Valid case with empty dotenv",
			dotenvContent: make(map[string]string),
			actualEnv: map[string]string{
				"API_URL":      "https://api.com/",
				"SECRET_TOKEN": "12345",
			},
			excpectedMap: make(map[string]string),
		},
		{
			name:      "Invalid case with missed variable in actual env without prefixes",
			keyPrefix: "VITE",
			dotenvContent: map[string]string{
				"VITE_API_URL":      "API_URL",
				"VITE_SECRET_TOKEN": "SECRET_TOKEN",
			},
			actualEnv: map[string]string{
				// missed variable
				// "API_URL":      "https://api.com/",
				"SECRET_TOKEN": "12345",
			},
			err: errMissedVariable,
		},
	}

	for _, tc := range testCases {
		tc := tc
		// hide logs
		log.Init(log.LogLevelDebug, true)

		t.Run(tc.name, func(t *testing.T) {
			prepareEnv(t, tc.actualEnv)

			entries := dotenv.ParseEntries(tc.dotenvContent, tc.keyPrefix, tc.placeholderPrefix)
			actualMap, err := mapPlaceholderToValue(entries)

			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, errMissedVariable.Error())
				return
			}

			if tc.err == nil && !reflect.DeepEqual(tc.excpectedMap, actualMap) {
				require.Failf(t, "Maps should be equal", "Maps aren't equal \nExpected: %#v \nActual: %#v", tc.excpectedMap, actualMap)
			}
		})
	}
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
