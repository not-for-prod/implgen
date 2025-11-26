package main

import (
	"os"

	"github.com/not-for-prod/implgen/generator"
	"github.com/not-for-prod/implgen/parser"
	"github.com/not-for-prod/implgen/pkg/clog"
	"github.com/not-for-prod/implgen/writer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// Flag names used in the CLI
	srcFlat                       = "src"
	dstFlat                       = "dst"
	interfaceNameFlag             = "interface-name"
	implementationNameFlag        = "impl-name"
	implementationPackageNameFlag = "impl-package"
	singleFileFlag                = "single-file"
	verboseFlag                   = "verbose"
)

const (
	// Default values for optional flags
	defaultImplementationName = "Implementation"
	defaultTracerName         = ""
)

// main defines and executes the CLI command using cobra.
// It parses input flags, runs code generation logic, and writes resulting files to disk.
func main() {
	// Define the root command
	cmd := &cobra.Command{
		Use:   "implgen",
		Short: "creates basic interface implementation",
		Long: `This tool generates Go implementations for interfaces in the given source package.
Example:
  implgen --src=./service --dst=./serviceimpl --interface-name=Greeter`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			parseCommand := flagsToParseCommand(flags)
			generateCommand := flagsToGenerateCommand(flags)
			writeCommand := flagsToWriteCommand(flags)

			_package, err := parseCommand.Execute()
			exitOnErr("failed to parse source file", err)

			// Run code generation using provided options
			files, err := generateCommand.Execute(_package)
			exitOnErr("failed to generate implementation", err)

			// Write generated files to disk
			err = writeCommand.Execute(files)
			exitOnErr("failed to write basic interfaces implementation", err)
		},
	}

	// Register command-line flags
	registerFlags(cmd)

	// Execute the root command
	err := cmd.Execute()
	exitOnErr("failed to execute command", err)
}

func exitOnErr(msg string, err error) {
	if err != nil {
		clog.Errorf("%s: %v", msg, err)
		os.Exit(1)
	}
}

// registerFlags - registers flags
func registerFlags(cmd *cobra.Command) {
	cmd.Flags().String(srcFlat, "", "source file")
	cmd.Flags().String(dstFlat, "", "destination file")
	_ = cmd.MarkFlagRequired(srcFlat)
	_ = cmd.MarkFlagRequired(dstFlat)

	cmd.Flags().String(interfaceNameFlag, "", "source interface name")
	cmd.Flags().Bool(singleFileFlag, false, "generate interface methods into single file")
	cmd.Flags().String(implementationNameFlag, defaultImplementationName, "generated implementation struct name")
	cmd.Flags().String(
		implementationPackageNameFlag, "",
		"generated implementation package name, can be used only when interface name is set",
	)
	cmd.Flags().Bool(verboseFlag, false, "enable verbose logging")
}

// flagsToParseCommand - parse cobra.Command flags into parser.Command
func flagsToParseCommand(flags *pflag.FlagSet) *parser.Command {
	// Parse required --src flag
	src, err := flags.GetString(srcFlat)
	if err != nil {
		clog.Errorf("failed to get src: %v", err)
		os.Exit(1)
	}

	return parser.NewCommand(src)
}

// flagsToGenerateCommand - parse cobra.Command flags into generator.Command
func flagsToGenerateCommand(flags *pflag.FlagSet) *generator.Command {
	// Parse required --dst flag
	dst, err := flags.GetString(dstFlat)
	if err != nil {
		clog.Errorf("failed to get dst: %v", err)
		os.Exit(1)
	}

	// Parse optional flags
	interfaceName, err := flags.GetString(interfaceNameFlag)
	if err != nil {
		clog.Errorf("failed to get interface-name: %v", err)
		os.Exit(1)
	}
	singleFile, _ := flags.GetBool(singleFileFlag)
	implementationName, _ := flags.GetString(implementationNameFlag)
	implementationPackageName, _ := flags.GetString(implementationPackageNameFlag)

	// Validate: impl package name requires interface name
	if implementationPackageName != "" && interfaceName == "" {
		clog.Errorf("flag %q requires %q to be set", implementationPackageNameFlag, interfaceNameFlag)
		os.Exit(1)
	}

	return generator.NewCommand(
		dst,
		interfaceName,             // src interface name
		implementationName,        // dst struct name
		implementationPackageName, // dst package name
		singleFile,
	)
}

func flagsToWriteCommand(flags *pflag.FlagSet) *writer.Command {
	verbose, _ := flags.GetBool(verboseFlag)

	return writer.NewCommand(verbose)
}
