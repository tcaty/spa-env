package shell

import (
	"fmt"
	"os"
	"os/user"

	"github.com/tcaty/spa-entrypoint/pkg/passwd"
)

// Find shell absolute path for current user
// return absolute path to shell and no error
// or return empty string and error if not found
func Find() (string, error) {
	shell := os.Getenv("SHELL")
	if shell != "" {
		return shell, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to get current user: %v", err)
	}

	entries, err := passwd.Parse()
	if err != nil {
		return "", fmt.Errorf("unable to parse passwd: %v", err)
	}

	for _, entry := range entries {
		if u.Uid == entry.Uid {
			return entry.Shell, nil
		}
	}

	return "", fmt.Errorf("shell for user %s not found", u.Username)
}
