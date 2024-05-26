package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcaty/spa-env/internal/common/log"
	"github.com/tcaty/spa-env/internal/generate"
)

type GenerateFlags struct {
	Workdir           string
	DotenvDev         string
	DotenvProd        string
	KeyPrefix         string
	PlaceholderPrefix string
	EnableComments    bool
	LogLevel          string
}

var generateFlags GenerateFlags

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Run generate command",
	Long:  "Generate .env file with placeholders for production mode based on development .env file",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := log.ValidateLogLevel(generateFlags.LogLevel); err != nil {
			return fmt.Errorf("--log-level validation failed: %v", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Init(generateFlags.LogLevel, false)

		log.Debug(
			"prefixes are initialized",
			"key_prefix", generateFlags.KeyPrefix,
			"placeholder_prefix", generateFlags.PlaceholderPrefix,
		)

		log.Info("starting generate .env file...")

		path, err := generate.Generate(
			generateFlags.Workdir,
			generateFlags.DotenvDev,
			generateFlags.DotenvProd,
			generateFlags.KeyPrefix,
			generateFlags.PlaceholderPrefix,
			generateFlags.EnableComments,
		)

		if err != nil {
			log.Fatal("error occured while generating .env file", err)
		}

		log.Info(
			".env file generated successfully",
			"path", path,
		)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&generateFlags.Workdir, "workdir", "w", "", "Path to working directory")
	generateCmd.PersistentFlags().StringVarP(&generateFlags.DotenvDev, "dotenv-dev", "", ".env.development", "Name of development .env file.")
	generateCmd.PersistentFlags().StringVarP(&generateFlags.DotenvProd, "dotenv-prod", "", ".env.production", "Name of production .env file.")
	generateCmd.PersistentFlags().StringVarP(&generateFlags.KeyPrefix, "key-prefix", "k", "", "Env variable prefix that will be parsed and generated")
	generateCmd.PersistentFlags().StringVarP(&generateFlags.PlaceholderPrefix, "placeholder-prefix", "p", "PLACEHOLDER", "Placeholder prefix that will be parsed and generated")
	generateCmd.PersistentFlags().BoolVarP(&generateFlags.EnableComments, "enable-comments", "", false, "Enable comments in generated .env file")
	generateCmd.PersistentFlags().StringVarP(&generateFlags.LogLevel, "log-level", "l", log.LogLevelInfo, "Log level")

	if err := generateCmd.MarkPersistentFlagRequired("workdir"); err != nil {
		return
	}
}
