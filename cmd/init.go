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

type InitFlags struct {
	WorkDir   string
	EnvFile   string
	EnvPrefix string
}

// TODO: rename to replace
var initFlags InitFlags

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(initFlags.WorkDir); err != nil {
			log.Fatalf("error occured while reading workdir: %v", err)
		}

		envFilePath, err := file.Find(initFlags.WorkDir, initFlags.EnvFile)
		if err != nil {
			log.Fatalf("error occured while finding .env file: %v", err)
		}

		envMap, err := env.MapEnvFileToActualEnv(envFilePath)
		if err != nil {
			log.Fatalf("error occured while reading .env file: %v", err)
		}

		err = filepath.WalkDir(initFlags.WorkDir, func(path string, d fs.DirEntry, err error) error {
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
	initCmd.PersistentFlags().StringVarP(&initFlags.WorkDir, "workdir", "", "", "Path to working directory")
	initCmd.PersistentFlags().StringVarP(&initFlags.EnvFile, "env-file", "", ".env", "Name of .env file")
	// TODO: replace to "use-prefix"
	initCmd.PersistentFlags().StringVarP(&initFlags.EnvPrefix, "env-prefix", "", "NEXT_PUBLIC_", "Prefix to env variables")

	initCmd.MarkPersistentFlagRequired("workdir")
}
