package command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/tcaty/spa-env/pkg/shell"
)

const (
	ShellForm = "shell"
	ExecForm  = "exec"
)

// Parse command string and return cmd in specified form
// Notice that shell will be found automatically for current user
// Example:
// cmd := "echo hello world"
// shell form: "/bin/bash -c 'echo hello world'"
// exec form: "echo hello world"
func Parse(cmd string, form string) (*exec.Cmd, error) {
	var name string
	var args []string

	switch form {
	case ShellForm:
		shell, err := shell.Find()
		if err != nil {
			return nil, fmt.Errorf("unable to find current user shell: %v", err)
		}
		name = shell
		args = []string{"-c", cmd}
	case ExecForm:
		arr := strings.Split(cmd, " ")
		name = arr[0]
		args = arr[1:]
	default:
		return nil, fmt.Errorf("unknown cmd form: %s", form)
	}

	command := exec.Command(name, args...)

	return command, nil
}

// Run passed command and print it output to stdout and stderr
func Run(cmd *exec.Cmd) error {
	cmdStdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create cmd stdout pipe: %v", err)
	}

	cmdStderrReader, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to create cmd stderr pipe: %v", err)
	}

	printOutput(cmdStdoutReader, os.Stdout)
	printOutput(cmdStderrReader, os.Stderr)

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("error occured while running cmd: %v", err)
	}

	return nil
}

// Print data to writer from reader
func printOutput(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			fmt.Fprintln(w, scanner.Text())
		}
	}()
}

func ValidateForm(form string) error {
	switch form {
	case ShellForm, ExecForm:
		return nil
	default:
		return fmt.Errorf(
			"form validation failed: expected %s or %s, but got %s",
			ShellForm, ExecForm, form,
		)
	}
}
