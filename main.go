package main

import (
	"github.com/not-for-prod/implgen/internal/implgen"
	"github.com/not-for-prod/implgen/pkg/logger"
	"github.com/not-for-prod/implgen/pkg/mockgen"
	"github.com/spf13/cobra"
)

const (
	srcFlag       = "src"
	dstFlag       = "dst"
	withOtelFlag  = "with-otel"
	interfaceName = "interface-name"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "implgen",
		Short: "creates interface implementation",
		Long:  `creates files for all interface methods`,
		Run: func(cmd *cobra.Command, _ []string) {
			src := cmd.Flag(srcFlag).Value.String()
			dst := cmd.Flag(dstFlag).Value.String()
			ifceName := cmd.Flag(interfaceName).Value.String()

			withOtel, err := cmd.Flags().GetBool(withOtelFlag)
			if err != nil {
				logger.Fatal(err.Error())
			}

			pkg, err := mockgen.SourceMode(src)
			if err != nil {
				logger.Fatalf(err.Error())
			}

			g := implgen.NewGenerator(pkg, src, dst, withOtel, ifceName)
			g.Generate()
		},
	}

	rootCmd.Flags().String(srcFlag, "", "path to interface")
	rootCmd.Flags().String(dstFlag, "", "path to generated files")
	rootCmd.Flags().Bool(withOtelFlag, false, "use otel tracer")
	rootCmd.Flags().String(interfaceName, "", "interface name")
	_ = rootCmd.Execute()
}
