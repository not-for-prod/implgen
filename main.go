package main

import (
	"os"

	"github.com/not-for-prod/implgen/generator"
	"github.com/not-for-prod/implgen/mockgen"
	"github.com/not-for-prod/implgen/mockgen/model"
	"github.com/not-for-prod/implgen/pkg/clog"
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
	enableTraceFlag               = "enable-trace"
	enableTestsFlag               = "enable-tests"
)

const (
	// Default values for optional flags
	defaultImplementationName = "Implementation"
)

// main defines and executes the CLI command using cobra.
// It parses input flags, runs code generation logic, and writes resulting files to disk.
func main() {
	// Define the root command
	cmd := &cobra.Command{
		Use:   "implgen",
		Short: "creates basic interface implementation",
		Long:  `This tool generates Go implementations for interfaces in the given source package.`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := cmd.Flags()
			pkg := parse(flags)
			generate(flags, pkg)
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
	cmd.Flags().Bool(
		enableTraceFlag, false,
		"whether you need to add otel.Tracer(<tracer-name>).Start(...) in your methods",
	)
	cmd.Flags().String(implementationNameFlag, defaultImplementationName, "generated implementation struct name")
	cmd.Flags().String(
		implementationPackageNameFlag, "",
		"generated implementation package name, can be used only when interface name is set",
	)
	cmd.Flags().Bool(enableTestsFlag, false, "generate interface method tests")
}

// parse - parse cobra.Command flags into mockgen.ParseCommand
func parse(flags *pflag.FlagSet) *model.Package {
	// Parse required --src flag
	src, err := flags.GetString(srcFlat)
	exitOnErr("failed to parse source file", err)

	pkg, err := mockgen.SourceMode(src)
	exitOnErr("failed to parse source file", err)

	return pkg
}

// generate - parse cobra.Command flags into generator.GenerateCommand
func generate(flags *pflag.FlagSet, pkg *model.Package) {
	// Parse required --src flag
	src, err := flags.GetString(srcFlat)
	exitOnErr("failed to parse source file", err)

	// Parse required --dst flag
	dst, err := flags.GetString(dstFlat)
	exitOnErr("failed to parse destination file", err)

	// Parse optional flags
	interfaceName, err := flags.GetString(interfaceNameFlag)
	exitOnErr("failed to parse interface name", err)

	implementationName, _ := flags.GetString(implementationNameFlag)
	implementationPackageName, _ := flags.GetString(implementationPackageNameFlag)
	enableTrace, _ := flags.GetBool(enableTraceFlag)
	enableTests, _ := flags.GetBool(enableTestsFlag)

	// Validate: impl package name requires interface name
	if implementationPackageName != "" && interfaceName == "" {
		clog.Errorf("flag %q requires %q to be set", implementationPackageNameFlag, interfaceNameFlag)
		os.Exit(1)
	}

	gc := generator.NewGenerateCommand(
		pkg,
		src,
		dst,
		interfaceName,
		implementationName,
		implementationPackageName,
		enableTrace,
		enableTests,
	)
	gc.Generate()
}
