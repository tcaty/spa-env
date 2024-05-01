package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/tcaty/spa-entrypoint/internal/replace"
	"github.com/tcaty/spa-entrypoint/pkg/command"
)

type ReplaceFlags struct {
	Workdir string
	Dotenv  string
	Cmd     string
	Form    string
	Verbose bool
}

var replaceFlags ReplaceFlags

var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace static env values from .env by values from actual env",
	Args: func(cmd *cobra.Command, args []string) error {
		switch replaceFlags.Form {
		case command.ShellForm, command.ExecForm:
			return nil
		default:
			return fmt.Errorf(
				"flags validation failed [--form]: wrong cmd form: expected %s or %s, but got %s",
				command.ShellForm, command.ExecForm, replaceFlags.Form,
			)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		log.Println("Starting environment variables replacement...")

		if err := replace.Replace(replaceFlags.Workdir, replaceFlags.Dotenv, replaceFlags.Verbose); err != nil {
			log.Fatalf("error while replacing: %v", err)
		}

		elapsed := time.Since(start)
		log.Printf("Replacement completed successfully in %s\n\n", elapsed)

		if replaceFlags.Cmd != "" {
			cmd, err := command.Parse(replaceFlags.Cmd, replaceFlags.Form)
			if err != nil {
				log.Fatalf("unable to parse cmd: %v", err)
			}

			if err := command.Run(cmd); err != nil {
				log.Fatalf("error occured while running cmd: %v", err)
			}
		}
	},
}

func init() {
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Workdir, "workdir", "w", "", "Path to working directory")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Dotenv, "dotenv", "d", ".env", "Name of .env file not path. It will be found automatically in workdir")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Cmd, "cmd", "c", "", "Command to execute after replacement")
	replaceCmd.PersistentFlags().StringVarP(&replaceFlags.Form, "form", "f", command.ExecForm, "Form in which command from --cmd will be run")
	replaceCmd.PersistentFlags().BoolVarP(&replaceFlags.Verbose, "verbose", "v", false, "Enable verbose logs")

	replaceCmd.MarkPersistentFlagRequired("workdir")
}
