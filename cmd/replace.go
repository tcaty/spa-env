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
)

type ReplaceFlags struct {
	WorkDir string
	EnvFile string
}

var replaceFlags ReplaceFlags

var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace static env values from .env by values from actual env",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(replaceFlags.WorkDir); err != nil {
			log.Fatalf("error occured while reading workdir: %v", err)
		}

		envFilePath, err := file.Find(replaceFlags.WorkDir, replaceFlags.EnvFile)
		if err != nil {
			log.Fatalf("error occured while finding .env file: %v", err)
		}

		envMap, err := env.MapEnvFileToActualEnv(envFilePath)
		if err != nil {
			log.Fatalf("error occured while mapping .env file to env: %v", err)
		}

		err = filepath.WalkDir(replaceFlags.WorkDir, func(path string, d fs.DirEntry, err error) error {
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

		// use args like shell commmand if they were passed
		if len(args) > 0 {
			shell := os.Getenv("SHELL")
			cmd := exec.Command(shell, "-c", strings.Join(args, " "))
			if err := cmd.Run(); err != nil {
				log.Fatalf("error occured while runnig: %v", err)
			}
		}
	},
}

func init() {
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.WorkDir, "workdir", "", "", "Path to working directory")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.EnvFile, "env-file", "", ".env", "Name of .env file")

	replaceCmd.MarkPersistentFlagRequired("workdir")
}
