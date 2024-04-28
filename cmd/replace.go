package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tcaty/spa-entrypoint/internal/env"
	"github.com/tcaty/spa-entrypoint/internal/file"
	"github.com/tcaty/spa-entrypoint/pkg/command"
	"github.com/tcaty/spa-entrypoint/pkg/shell"
)

type ReplaceFlags struct {
	Workdir string
	Dotenv  string
	Cmd     string
}

var replaceFlags ReplaceFlags

var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace static env values from .env by values from actual env",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(replaceFlags.Workdir); err != nil {
			log.Fatalf("error occured while reading workdir: %v", err)
		}

		dotenvPath, err := file.Find(replaceFlags.Workdir, replaceFlags.Dotenv)
		if err != nil {
			log.Fatalf("error occured while finding .env file: %v", err)
		}

		envMap, err := env.MapDotenvToActualEnv(dotenvPath)
		if err != nil {
			log.Fatalf("error occured while mapping .env file to env: %v", err)
		}

		err = filepath.WalkDir(replaceFlags.Workdir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("prevent panic by handling failure accessing a path %q: %v", path, err)
			}

			// skip unnecessary files
			if strings.Contains(path, "node_modules") || filepath.Ext(path) != ".js" {
				return nil
			}

			if err := file.Replace(path, envMap); err != nil {
				return fmt.Errorf("error occured while replacing file content: %v", err)
			}

			return nil
		})

		if err != nil {
			log.Fatalf("error occured in main logic: %v", err)
		}

		// exec shell command if it was passed
		if replaceFlags.Cmd != "" {
			shell, err := shell.Find()
			if err != nil {
				log.Fatalf("unable to find current user shell: %v", err)
			}

			cmd := exec.Command(shell, "-c", replaceFlags.Cmd)

			if err := command.Run(cmd); err != nil {
				log.Fatalf("error occured while running cmd: %v", err)
			}
		}
	},
}

func init() {
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Workdir, "workdir", "", "", "Path to working directory")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Dotenv, "dotenv", "", ".env", "Name of .env file")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Cmd, "cmd", "", "", "Command to execute after replacement")

	replaceCmd.MarkPersistentFlagRequired("workdir")
}
