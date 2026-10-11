package parser_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/not-for-prod/implgen/generator"
	"github.com/not-for-prod/implgen/parser"
	"github.com/not-for-prod/implgen/writer"
)

func TestGenerateSourcePackageTypes(t *testing.T) {
	for _, tc := range []struct {
		name       string
		singleFile bool
		relative   bool
	}{
		{name: "multiple files, absolute source"},
		{name: "single file, relative source", singleFile: true, relative: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			fixtures := map[string]string{
				"go.mod": "module example.com/fixture\n\ngo 1.25.0\n",
				"source/interface.go": `package source

import (
	"context"
	"example.com/fixture/dto"
)

type ERequest struct{}
type EResponse struct{}

type Interface interface {
	D(ctx context.Context, req int, opts ...dto.Request) error
	E(ctx context.Context, req ERequest) (EResponse, error)
	F(ctx context.Context, req FRequest) (FResponse, error)
}
`,
				"source/types.go": "package source\n\ntype FRequest struct{}\ntype FResponse struct{}\n",
				"dto/types.go":    "package dto\n\ntype Request struct{}\n",
			}
			for name, content := range fixtures {
				writeFixture(t, filepath.Join(root, name), content)
			}

			src := filepath.Join(root, "source/interface.go")
			if tc.relative {
				cwd, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				src, err = filepath.Rel(cwd, src)
				if err != nil {
					t.Fatal(err)
				}
			}
			pkg, err := parser.NewCommand(src).Execute()
			if err != nil {
				t.Fatal(err)
			}
			if got := pkg.Imports[0]; got.Alias != "source" || got.Path != "example.com/fixture/source" {
				t.Fatalf("source import = %+v, want source example.com/fixture/source", got)
			}

			files, err := generator.NewCommand(filepath.Join(root, "out"), "Interface", "Test", "test", tc.singleFile).Execute(pkg)
			if err != nil {
				t.Fatal(err)
			}
			if err := writer.NewCommand(false).Execute(files); err != nil {
				t.Fatal(err)
			}
			writeFixture(t, filepath.Join(root, "out/test/implements_test.go"), `package test

import "example.com/fixture/source"

var _ source.Interface = (*Test)(nil)
`)
			cmd := exec.Command("go", "test", "./...")
			cmd.Dir = root
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("generated implementation does not compile: %v\n%s", err, output)
			}
		})
	}
}

func writeFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
