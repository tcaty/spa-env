package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tcaty/spa-env/internal/common/log"
	"github.com/tcaty/spa-env/internal/replace"
	"github.com/tcaty/spa-env/pkg/command"
)

type ReplaceFlags struct {
	Workdir           string
	Dotenv            string
	KeyPrefix         string
	PlaceholderPrefix string
	Cmd               string
	CmdForm           string
	LogLevel          string
}

var replaceFlags ReplaceFlags

var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Run replace command",
	Long:  "This commmand replaces static env values from .env by values from actual environment",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := command.ValidateForm(replaceFlags.CmdForm); err != nil {
			return fmt.Errorf("--form validation failed: %v", err)
		}
		if err := log.ValidateLogLevel(replaceFlags.LogLevel); err != nil {
			return fmt.Errorf("--log-level validation failed: %v", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Init(replaceFlags.LogLevel, false)

		start := time.Now()
		log.Info("Starting environment variables replacement...")

		filesUpdated, err := replace.Replace(
			replaceFlags.Workdir,
			replaceFlags.Dotenv,
			replaceFlags.KeyPrefix,
			replaceFlags.PlaceholderPrefix,
		)

		if err != nil {
			log.Fatal("error occured while replacing", err)
		}

		duration := time.Since(start)
		log.Info(
			"replacement completed successfully",
			"duration", duration,
			"files_updated", filesUpdated,
		)

		if replaceFlags.Cmd != "" {
			cmd, err := command.Parse(replaceFlags.Cmd, replaceFlags.CmdForm)
			if err != nil {
				log.Fatal("unable to parse cmd", err)
			}

			if err := command.Run(cmd); err != nil {
				log.Fatal("error occured while running cmd", err)
			}
		}
	},
}

func init() {
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Workdir, "workdir", "w", "", "Path to working directory")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Dotenv, "dotenv", "d", ".env", "Name of .env file not path. It will be found automatically in workdir")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.KeyPrefix, "key-prefix", "k", "", "Env variable prefix that will be parsed and replaced")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.PlaceholderPrefix, "placeholder-prefix", "p", "PLACEHOLDER", "Placeholder prefix that will be parsed and replaced")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Cmd, "cmd", "c", "", "Command to execute after replacement")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.CmdForm, "cmd-form", "f", command.ExecForm, "Form in which command from --cmd will be run")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.LogLevel, "log-level", "l", log.LogLevelInfo, "Log level")

	if err := replaceCmd.MarkPersistentFlagRequired("workdir"); err != nil {
		return
	}
}
