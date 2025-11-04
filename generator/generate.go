package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/not-for-prod/implgen/model"
	stringCase "github.com/not-for-prod/implgen/pkg/string-case"
	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateCommand holds configuration for generating a Go implementation
// of an interface, including destination, naming, and output structure.
type GenerateCommand struct {
	// dst is the base directory where generated files will be written.
	dst string

	// interfaceName is the name of the interface to generate.
	// If empty, all interfaces in the package will be processed.
	interfaceName string

	// implementationName is the name of the struct that implements the interface.
	implementationName string

	// implementationPackageName overrides the default package name for the generated code.
	// If empty, the package name is derived from the interface name.
	implementationPackageName string

	// singleFile determines whether all methods should be generated into a single file.
	// If false, each method will be written into its own file.
	singleFile bool
}

// NewGenerateCommand creates a new GenerateCommand with the given parameters.
// This struct is typically passed into a code generator to drive its behavior.
func NewGenerateCommand(
	dst string,
	interfaceName string,
	implementationName string,
	implementationPackageName string,
	singleFile bool,
) *GenerateCommand {
	return &GenerateCommand{
		dst:                       dst,
		interfaceName:             interfaceName,
		implementationName:        implementationName,
		implementationPackageName: implementationPackageName,
		singleFile:                singleFile,
	}
}

// packageName determines the Go package name for the generated code
// based on the interface name and optional overrides provided in the generator.
func (cmd *GenerateCommand) packageName(ifce model.Interface) string {
	packageName := ifce.Name

	if cmd.interfaceName != "" && cmd.implementationPackageName != "" {
		packageName = cmd.implementationPackageName
	}

	return stringCase.SnakeCase(packageName)
}

// folderName determines the Go package folder name for the generated code
// based on the interface name and optional overrides provided in the generator.
func (cmd *GenerateCommand) folderName(ifce model.Interface) string {
	pkgName := cmd.packageName(ifce)

	return strings.Replace(pkgName, "_", "-", -1)
}

// dstPath builds the output directory path for a given interface's implementation,
// using the base destination directory and the generated package name.
func (cmd *GenerateCommand) dstPath(ifce model.Interface) string {
	return filepath.Join(cmd.dst, cmd.folderName(ifce))
}

// Execute creates basic implementations for all interfaces in the provided package.
// If a specific interface name is configured, only that one is processed.
func (cmd *GenerateCommand) Execute(pkg model.Package) ([]model.File, error) {
	files := make([]model.File, 0)

	for _, _interface := range pkg.Interfaces {
		if cmd.interfaceName == "" || cmd.interfaceName == _interface.Name {
			file, err := cmd.generateInterface(pkg, _interface)
			if err != nil {
				return nil, fmt.Errorf("failed to generate interface for %s, err: %w", _interface.Name, err)
			}
			files = append(files, file...)
		}
	}

	return files, nil
}

// generateInterface generates a full implementation of the given interface,
// including its struct declaration, constructor, and method stubs.
// Methods are split into files if configured via i.singleFile.
func (cmd *GenerateCommand) generateInterface(pkg model.Package, ifce model.Interface) ([]model.File, error) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile("", "")
	files := make([]model.File, 0)

	cmd.generateHeader(g, pkg, ifce)

	g.P("type ", cmd.implementationName, " struct {")
	g.P("}")
	g.P()
	g.P("func New", cmd.implementationName, "() *", cmd.implementationName, " {")
	g.P("return &", cmd.implementationName, "{}")
	g.P("}")

	for _, method := range ifce.Methods {
		if cmd.singleFile {
			g.P()
			cmd.generateMethod(g, method)
		} else {
			file, err := cmd.generateMethodFile(pkg, ifce, method)
			if err != nil {
				return nil, err
			}

			files = append(files, file)
		}
	}

	content, err := g.Content()
	if err != nil {
		return nil, err
	}

	files = append(
		files, model.File{
			Path: filepath.Join(cmd.dstPath(ifce), stringCase.KebabCase(cmd.implementationName)+".go"),
			Data: content,
		},
	)

	return files, nil
}

// generateHeader writes the file header including package declaration and imports.
func (cmd *GenerateCommand) generateHeader(g *protogen.GeneratedFile, pkg model.Package, ifce model.Interface) {
	g.P("package ", cmd.packageName(ifce))
	g.P()
	generateImports(g, pkg)
	g.P()
}

// generateImports writes import statements for a given package,
// including all user-defined imports and a default one for OpenTelemetry.
func generateImports(g *protogen.GeneratedFile, pkg model.Package) {
	g.P("import (")
	for _, _import := range pkg.Imports {
		g.P(_import.Alias, " \"", _import.Path, "\"")
	}
	g.P(")")
}

// generateMethodFile generates a standalone Go file containing a single method implementation stub.
func (cmd *GenerateCommand) generateMethodFile(
	pkg model.Package,
	ifce model.Interface,
	method model.Method,
) (model.File, error) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile("", "")

	cmd.generateHeader(g, pkg, ifce)
	cmd.generateMethod(g, method)

	content, err := g.Content()
	if err != nil {
		return model.File{}, fmt.Errorf(
			"failed to generate method for %s.%s, err: %w",
			ifce.Name,
			method.Name,
			err,
		)
	}

	return model.File{
		Path: filepath.Join(cmd.dstPath(ifce), stringCase.SnakeCase(method.Name)+".go"),
		Data: content,
	}, nil
}

// generateMethod writes the method implementation stub to the provided generated file.
func (cmd *GenerateCommand) generateMethod(g *protogen.GeneratedFile, method model.Method) {
	params := generateParams(method.In)
	results := generateResults(method.Out)

	g.P("func (i *", cmd.implementationName, ")", method.Name, " ", params, " ", results, "{")
	g.P("panic(\"implement me\")")
	g.P("}")
	g.P()
}

// generateParams builds a function parameter list from a slice of Parameter structs.
func generateParams(params []model.Parameter) string {
	b := strings.Builder{}
	b.WriteString("(")

	for i, param := range params {
		if param.Type == "context.Context" {
			b.WriteString("ctx")
		} else {
			b.WriteString(param.Name)
		}

		b.WriteString(" ")
		b.WriteString(param.Type)
		if i != len(params)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString(")")
	return b.String()
}

// generateResults builds a function result list from a slice of Parameter structs.
// It wraps the result list in parentheses only if there is more than one result.
func generateResults(results []model.Parameter) string {
	b := strings.Builder{}

	if len(results) > 1 {
		b.WriteString("(")
	}

	for i, result := range results {
		b.WriteString(result.Type)
		if i != len(results)-1 {
			b.WriteString(", ")
		}
	}

	if len(results) > 1 {
		b.WriteString(")")
	}

	return b.String()
}
