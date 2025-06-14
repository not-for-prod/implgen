package repo

import (
	"github.com/not-for-prod/implgen/pkg/clog"
	"github.com/not-for-prod/implgen/pkg/mockgen"
	"github.com/spf13/cobra"
)

const (
	srcFlag           = "src"
	dstFlag           = "dst"
	interfaceNameFlag = "interface-name"
	withTestsFlag     = "tests"
)

func NewCMD() *cobra.Command {
	c := &cobra.Command{
		Use:   "repo",
		Short: "creates repository interface implementation",
		Long:  `creates files for all interface methods with db layer stuff`,
		Run: func(cmd *cobra.Command, _ []string) {
			src, err := cmd.Flags().GetString(srcFlag)
			if err != nil {
				clog.Fatal(err.Error())
			}

			dst, err := cmd.Flags().GetString(dstFlag)
			if err != nil {
				clog.Fatal(err.Error())
			}

			ifceName, err := cmd.Flags().GetString(interfaceNameFlag)
			if err != nil {
				clog.Fatal(err.Error())
			}

			pkg, err := mockgen.SourceMode(src)
			if err != nil {
				clog.Fatal(err.Error())
			}

			withTests, err := cmd.Flags().GetBool(withTestsFlag)
			if err != nil {
				clog.Fatal(err.Error())
			}

			newRepoGenerator(src, dst, ifceName, pkg, withTests).generate()
		},
	}

	c.Flags().String(srcFlag, "", "path to interface")
	c.Flags().String(dstFlag, "", "path to generated files")
	c.Flags().String(interfaceNameFlag, "", "interface name")
	c.Flags().Bool(withTestsFlag, false, "generate test files")

	return c
}
