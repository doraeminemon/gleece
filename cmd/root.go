package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/gopher-fleece/gleece/cmd/arguments"
	"github.com/gopher-fleece/gleece/infrastructure/logger"
	"github.com/spf13/cobra"
)

var cliArgs = arguments.CliArguments{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gleece",
	Short: "Gleece - A Simplified Framework for Building REST APIs in Go",
	Long: fmt.Sprintf(
		"\n"+
			"Gleece - A Simplified Framework for Building REST APIs in Go"+
			"\n\n\n"+
			"%s\n\n"+
			"%s\n\n",
		arguments.GopherAscii,
		arguments.GleeceIntroDoc,
	),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !cliArgs.NoBanner {
			fmt.Println(arguments.GopherAscii)
		}

		// Handle the verbosity flag here if you want it executed for every subcommand
		if cmd.Flag("verbosity") == nil {
			logger.SetLogLevel(logger.LogLevelInfo)
			return
		}

		verbosityInt, err := cmd.Flags().GetUint8("verbosity")
		if err != nil {
			logger.SetLogLevel(logger.LogLevelAll)
			logger.Warn("Could not obtain verbosity level from arguments. Fell back to 'all'. Error - %v", err)
			return
		}

		verbosity := logger.LogLevel(verbosityInt)
		logger.SetLogLevel(verbosity)
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info(`Gleece called with no parameters. Assuming 'generate spec-and-routes -c "./gleece.config.json"'`)
		err := GenerateSpecAndRoutes(arguments.CliArguments{ConfigPath: "./gleece.config.json"})
		if err != nil {
			logger.Fatal("Failed to generate spec and routes: %v", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// ExecuteWithArgs runs the root command with provided arguments and captures stdout, stderr, and logs.
func ExecuteWithArgs(args []string, redirectLogs bool) arguments.ExecuteWithArgsResult {
	// Capture the original log output
	originalLogOutput := log.Writer()

	// Create buffers to capture the command output and logs
	var out bytes.Buffer
	var errBuf bytes.Buffer
	var logBuf bytes.Buffer

	// Redirect log output to the logBuf temporarily
	log.SetOutput(&logBuf)

	// Set the command's output streams
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&errBuf)
	rootCmd.SetArgs(args)

	// Defer restoring the log output and capture result
	defer func() {
		// Restore original log output
		log.SetOutput(originalLogOutput)
	}()

	// Execute the root command
	err := rootCmd.Execute()

	// Return the captured result
	return arguments.ExecuteWithArgsResult{
		Error:  err,
		StdOut: out.String(),
		StdErr: errBuf.String(),
		Logs:   logBuf.String(), // Include the captured logs
	}
}

func init() {
	initGenerateCommandHierarchy()
	rootCmd.PersistentFlags().BoolVar(
		&cliArgs.NoBanner,
		"no-banner",
		false,
		"Disables the CLI banner",
	)

	rootCmd.PersistentFlags().Uint8VarP(
		&cliArgs.Verbosity,
		"verbosity",
		"v",
		2,
		"Determines the verbosity of Gleece's traces. 0 = Output everything, 5 = Output fatal errors only",
	)

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(generateCmd)
}
