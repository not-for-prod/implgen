package main

import (
	"github.com/not-for-prod/implgen/internal/implgen"
	"github.com/not-for-prod/implgen/internal/mockgen"
	"github.com/not-for-prod/implgen/internal/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	srcFlag      = "src"
	dstFlag      = "dst"
	withOtelFlag = "with-otel"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "implgen",
		Short: "creates interface implementation",
		Long:  `creates files for all interface methods`,
		Run: func(cmd *cobra.Command, _ []string) {
			src := cmd.Flag(srcFlag).Value.String()
			dst := cmd.Flag(dstFlag).Value.String()

			withOtel, err := cmd.Flags().GetBool(withOtelFlag)
			if err != nil {
				logger.Fatal(err.Error())
			}

			pkg, err := mockgen.SourceMode(src)
			if err != nil {
				logger.Fatalf(err.Error())
			}

			g := implgen.NewGenerator(pkg, src, dst, withOtel)
			g.Generate()
		},
	}

	rootCmd.Flags().String(srcFlag, "", "path to interface")
	rootCmd.Flags().String(dstFlag, "", "path to generated files")
	rootCmd.Flags().Bool(withOtelFlag, false, "use otel tracer")
	_ = rootCmd.Execute()
}
