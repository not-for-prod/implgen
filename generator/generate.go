package generator

import (
	"fmt"
	"go/token"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/not-for-prod/implgen/mockgen"
	"github.com/not-for-prod/implgen/mockgen/model"
	"github.com/not-for-prod/implgen/pkg/clog"
	"github.com/not-for-prod/implgen/pkg/strtools"
	"github.com/not-for-prod/implgen/writer"
	"github.com/samber/lo"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	otelPackage    = "go.opentelemetry.io/otel"
	testingPackage = "testing"
	testifyPackage = "github.com/stretchr/testify/suite"
)

// GenerateCommand holds configuration for generating a Go implementation
// of an interface, including destination, naming, and output structure.
type GenerateCommand struct {
	packageMap map[string]string
	pkg        *model.Package
	src        string

	// dst is the base directory where generated files will be written.
	dst string

	// interfaceName is the name of the interface to Generate.
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

	// enableTrace determines whether otel.Tracer(tracerName).Start(...) will be used inside method.
	enableTrace bool

	// tracerName is the name used in otel.Tracer(tracerName).Start(...) for tracing spans.
	tracerName string

	enableTests bool
}

func NewGenerateCommand(
	pkg *model.Package,
	src, dst string,
	interfaceName string,
	implementationName string,
	implementationPackage string,
	singleFile bool,
	enableTests bool,
) *GenerateCommand {
	g := &GenerateCommand{
		pkg:                       pkg,
		src:                       src,
		dst:                       dst,
		interfaceName:             interfaceName,
		implementationName:        implementationName,
		implementationPackageName: implementationPackage,
		singleFile:                singleFile,
		enableTests:               enableTests,
	}

	// Get all required imports, and Generate unique names for them all.
	im := pkg.Imports()

	// Sort keys to make import alias generation predictable
	sortedPaths := make([]string, len(im))
	x := 0
	for pth := range im {
		sortedPaths[x] = pth
		x++
	}
	sort.Strings(sortedPaths)

	packagesName := mockgen.CreatePackageMap(sortedPaths)

	g.packageMap = make(map[string]string, len(im))
	localNames := make(map[string]bool, len(im))
	for _, pth := range sortedPaths {
		base, ok := packagesName[pth]
		if !ok {
			base = mockgen.Sanitize(path.Base(pth))
		}

		// Local names for an imported package can usually be the basename of the import path.
		// A couple of situations don't permit that, such as duplicate local names
		// (e.g. importing "html/template" and "text/template"), or where the basename is
		// a keyword (e.g. "foo/case").
		// try base0, base1, ...
		pkgName := base
		i := 0
		for localNames[pkgName] || token.Lookup(pkgName).IsKeyword() {
			pkgName = base + strconv.Itoa(i)
			i++
		}

		g.packageMap[pth] = pkgName
		localNames[pkgName] = true
	}

	return g
}

func (gc *GenerateCommand) Generate() {
	for _, ifce := range gc.pkg.Interfaces {
		if gc.interfaceName == "" || gc.interfaceName == ifce.Name {
			if gc.enableTests {
				gc.generateImplementationTest(ifce)
			}
			gc.generateImplementation(ifce)
		}
	}
}

// packageName determines the Go package name for the generated code
// based on the interface name and optional overrides provided in the generator.
func (gc *GenerateCommand) packageName(ifce *model.Interface) string {
	packageName := ifce.Name

	if gc.interfaceName != "" && gc.implementationPackageName != "" {
		packageName = gc.implementationPackageName
	}

	return strtools.GoPackageCase(packageName)
}

// dstPath builds the output directory path for a given interface's implementation,
// using the base destination directory and the generated package name.
func (gc *GenerateCommand) dstPath(ifce *model.Interface) string {
	return filepath.Join(gc.dst, gc.packageName(ifce))
}

func (gc *GenerateCommand) getArgNames(m *model.Method) []string {
	argNames := make([]string, len(m.In))
	for i, p := range m.In {
		name := p.Name
		if name == "" || name == "_" {
			name = fmt.Sprintf("arg%d", i)
		}
		argNames[i] = name
	}
	if m.Variadic != nil {
		name := m.Variadic.Name
		if name == "" {
			name = fmt.Sprintf("arg%d", len(m.In))
		}
		argNames = append(argNames, name) //nolint:makezero // mockgen authors
	}
	return argNames
}

func (gc *GenerateCommand) getArgTypes(m *model.Method, pkgOverride string) []string {
	argTypes := make([]string, len(m.In))
	for i, p := range m.In {
		argTypes[i] = p.Type.String(gc.packageMap, pkgOverride)
	}
	if m.Variadic != nil {
		//nolint:makezero // mockgen authors
		argTypes = append(
			argTypes,
			"..."+m.Variadic.Type.String(gc.packageMap, pkgOverride),
		)
	}
	return argTypes
}

func (gc *GenerateCommand) generateHeader(g *protogen.GeneratedFile, ifce *model.Interface) {
	g.P(
		"package ",
		strtools.SnakeCase(lo.Ternary(gc.implementationPackageName == "", ifce.Name, gc.implementationPackageName)),
	)
	g.P()
	g.P("import (")
	for imp := range gc.packageMap {
		g.P("\t\"", imp, "\"") // cannot use g.Import here
	}
	g.P("\t\"", testingPackage, "\"")
	g.P("\t\"", testifyPackage, "\"")
	g.P("\t\"", otelPackage, "\"")
	g.P(")")
	g.P()
}

func (gc *GenerateCommand) generateImplementation(ifce *model.Interface) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile("", protogen.GoImportPath(gc.dst))
	path := gc.dstPath(ifce) + "/implementation.go"

	gc.generateHeader(g, ifce)
	g.P("type Implementation struct {")
	g.P("}")
	g.P()
	g.P("func NewImplementation() *Implementation {")
	g.P("\treturn &Implementation{}")
	g.P("}")

	err := writer.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}

	for _, method := range ifce.Methods {
		gc.generateMethod(ifce, method)
	}
}

func (gc *GenerateCommand) generateImplementationTest(ifce *model.Interface) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile("", protogen.GoImportPath(gc.dst))
	path := gc.dstPath(ifce) + "/implementation_test.go"

	gc.generateHeader(g, ifce)
	g.P("type TestSuite struct{")
	g.P("suite.Suite")
	g.P("impl *Implementation")
	g.P("}")
	g.P()
	g.P("func (suite *TestSuite) SetupSuite() {")
	g.P("}")
	g.P()
	g.P("func (suite *TestSuite) TearDownSuite() {")
	g.P("}")
	g.P()
	g.P("func TestTestSuite(t *testing.T) {")
	g.P("suite.Run(t, new(TestSuite))")
	g.P("}")

	err := writer.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}

	for _, method := range ifce.Methods {
		gc.generateMethodTest(ifce, method)
	}
}

func (gc *GenerateCommand) generateMethod(ifce *model.Interface, m *model.Method) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile("", protogen.GoImportPath(gc.dst))
	path := gc.dstPath(ifce) + "/" + strtools.SnakeCase(m.Name) + ".go"

	gc.generateHeader(g, ifce)

	argNames := gc.getArgNames(m)
	argTypes := gc.getArgTypes(m, "")
	argString := mockgen.MakeArgString(argNames, argTypes)

	rets := make([]string, len(m.Out))
	for i, p := range m.Out {
		rets[i] = p.Type.String(gc.packageMap, "")
	}
	retString := strings.Join(rets, ", ")
	if len(rets) > 1 {
		retString = "(" + retString + ")"
	}
	if retString != "" {
		retString = " " + retString
	}

	g.P(
		fmt.Sprintf(
			"func (i *Implementation) %v(%v)%v {", m.Name, argString,
			retString,
		),
	)

	if len(argTypes) > 0 && strings.HasPrefix(argTypes[0], "context.") && gc.enableTrace {
		g.P(
			argNames[0], ", span := otel.Tracer(\"\").Start(", argNames[0], ", \"", ifce.Name, "Implementation.",
			m.Name, "\")",
		)
		g.P("defer span.End()")
		g.P()
	}

	g.P("\tpanic(\"implement me\")")
	g.P("}")

	err := writer.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}
}

func (gc *GenerateCommand) generateMethodTest(ifce *model.Interface, m *model.Method) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile("", protogen.GoImportPath(gc.dst))
	path := gc.dstPath(ifce) + "/" + strtools.SnakeCase(m.Name) + "_test.go"

	gc.generateHeader(g, ifce)
	g.P("func (suite *TestSuite) Test", m.Name, "() {")
	g.P("}")

	err := writer.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}
}
