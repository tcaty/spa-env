package command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

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
