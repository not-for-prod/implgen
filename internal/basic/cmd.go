package basic

import (
	"github.com/not-for-prod/implgen/pkg/clog"
	"github.com/not-for-prod/implgen/pkg/mockgen"
	"github.com/spf13/cobra"
)

const (
	srcFlag       = "src"
	dstFlag       = "dst"
	interfaceName = "interface-name"
	withTestsFlag = "tests"
)

func NewCMD() *cobra.Command {
	c := &cobra.Command{
		Use:   "basic",
		Short: "creates interface implementation",
		Long:  `creates files for all interface methods`,
		Run: func(cmd *cobra.Command, _ []string) {
			src := cmd.Flag(srcFlag).Value.String()
			dst := cmd.Flag(dstFlag).Value.String()
			ifceName := cmd.Flag(interfaceName).Value.String()

			withTests, err := cmd.Flags().GetBool(withTestsFlag)
			if err != nil {
				clog.Fatal(err.Error())
			}

			pkg, err := mockgen.SourceMode(src)
			if err != nil {
				clog.Fatal(err.Error())
			}

			g := newGenerator(pkg, src, dst, ifceName, withTests)
			g.generate()
		},
	}

	c.Flags().String(srcFlag, "", "path to interface")
	c.Flags().String(dstFlag, "", "path to generated files")
	c.Flags().String(interfaceName, "", "interface name")
	c.Flags().Bool(withTestsFlag, false, "generate test files")

	return c
}
