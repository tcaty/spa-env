package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
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

// TODO: rename to replace
var replaceFlags ReplaceFlags

// replaceCmd represents the replace command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		fmt.Println(envFilePath)

		// TODO: run command from args
	},
}

func init() {
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.WorkDir, "workdir", "", "", "Path to working directory")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.EnvFile, "env-file", "", ".env", "Name of .env file")

	replaceCmd.MarkPersistentFlagRequired("workdir")
}
