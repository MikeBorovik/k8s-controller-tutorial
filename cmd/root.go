/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	logLevel  string
	logFormat string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-controller-tutorial",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to k8s-controller-tutorial CLI!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Logger flags
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level: trace, debug, info, warn, error, none")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "json", "Log format: json, console")
	
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(initializeLogger)
}

func initializeLogger() {
	level := parseLogLevel(logLevel)
	zerolog.SetGlobalLevel(level)

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"

	baseLogger := zerolog.New(os.Stderr).With().Timestamp()
	
	// caller for trace level
	if level == zerolog.TraceLevel {
		baseLogger = baseLogger.Caller()
	}

	if strings.ToLower(logFormat) == "console" {
		output := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
		}
		log.Logger = baseLogger.Logger().Output(output)
	} else {
		log.Logger = baseLogger.Logger()
	}
	
	log.Debug().Str("format", logFormat).Str("level", logLevel).
		Msg("Logger initialized")
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "none":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}
