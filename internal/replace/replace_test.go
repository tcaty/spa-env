package replace

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
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
				"VITE_API_URL":      "API_URL",
				"VITE_SECRET_TOKEN": "SECRET_TOKEN",
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
			keyPrefix: "NEXT_PUBLIC",
			dotenvContent: map[string]string{
				"POSTGRES_CONN_STRING": "POSTGRES_CONN_STRING",
				"NEXT_PUBLIC_API_URL":  "API_URL",
				"NEXT_PUBLIC_TOKEN":    "SECRET_TOKEN",
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
			placeholderPrefix: "PLACEHOLDER",
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
			keyPrefix:         "NEXT_PUBLIC",
			placeholderPrefix: "PLACEHOLDER",
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
			name: "Invalid case with missed variable in actual env without prefixes",
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
		{
			name: "Invalid case with missed variable in actual env caused by missed placeholderPrefix",
			// missed placeholder prefix
			// placeholderPrefix: "PLACEHOLDER",
			dotenvContent: map[string]string{
				"VITE_API_URL":      "PLACEHOLDER_API_URL",
				"VITE_SECRET_TOKEN": "PLACEHOLDER_SECRET_TOKEN",
			},
			actualEnv: map[string]string{
				"API_URL":      "https://api.com/",
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

			actualMap, err := mapPlaceholderToValue(tc.dotenvContent, tc.keyPrefix, tc.placeholderPrefix)

			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, errMissedVariable.Error())
			}

			if tc.err == nil && !reflect.DeepEqual(tc.excpectedMap, actualMap) {
				require.Failf(t, "Maps should be equal", "Maps aren't equal \nExpected: %#v \nActual: %#v", tc.excpectedMap, actualMap)
			}
		})
	}
}

func TestGetenv(t *testing.T) {
	tcs := []struct {
		name      string
		key       string
		prefix    string
		actualEnv map[string]string
		expected  string
	}{
		{
			name: "Without prefix",
			key:  "API_URL",
			actualEnv: map[string]string{
				"API_URL":      "https://api.com/",
				"SECRET_TOKEN": "12345",
			},
			expected: "https://api.com/",
		},
		{
			name:   "With right prefix without suffix",
			key:    "PLACEHOLDER_API_URL",
			prefix: "PLACEHOLDER",
			actualEnv: map[string]string{
				"PLACEHOLDER_API_URL": "https://wrong.api.com/",
				"API_URL":             "https://api.com/",
			},
			expected: "https://api.com/",
		},
		{
			name:   "With right prefix with suffix",
			key:    "PLACEHOLDER_API_URL",
			prefix: "PLACEHOLDER_",
			actualEnv: map[string]string{
				"PLACEHOLDER_API_URL": "https://wrong.api.com/",
				"API_URL":             "https://api.com/",
			},
			expected: "https://api.com/",
		},
		{
			name: "Without prefix but key miseed in environment",
			key:  "API_URL",
			actualEnv: map[string]string{
				"SECRET_TOKEN": "12345",
			},
			expected: "",
		},
		{
			name: "With wrong prefix",
			// It is assumed that we want to get API_URL from environment
			// but there will be nothing deleted from key because it doesn't contain prefix
			// therefore PLACEHOLDER_API_URL will be returned
			key:    "PLACEHOLDER_API_URL",
			prefix: "WRONGPREFIX",
			actualEnv: map[string]string{
				"PLACEHOLDER_API_URL": "https://wrong.api.com/",
				"API_URL":             "https://api.com/",
			},
			expected: "https://wrong.api.com/",
		},
	}

	for _, tc := range tcs {
		tc := tc
		// hide logs
		log.Init(log.LogLevelDebug, true)

		t.Run(tc.name, func(t *testing.T) {
			prepareEnv(t, tc.actualEnv)

			value := getenv(tc.key, tc.prefix)

			require.Equal(t, tc.expected, value)
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
