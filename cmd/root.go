package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spa-env",
	Short: "Run spa-env",
	Long:  "spa-env is a set of different useful utils that helps to work with environment variables in spa",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// TODO: generate - generate .env file for production based on .env.development
	// TODO: validate - validate code for environment variables usage
	rootCmd.AddCommand(replaceCmd)
}
