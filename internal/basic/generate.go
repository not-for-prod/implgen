package basic

import (
	"fmt"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/not-for-prod/implgen/pkg/clog"
	"github.com/not-for-prod/implgen/pkg/fwriter"
	"github.com/not-for-prod/implgen/pkg/mockgen"
	"github.com/not-for-prod/implgen/pkg/mockgen/model"
	"github.com/not-for-prod/implgen/pkg/strtools"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	otelPackage    = "go.opentelemetry.io/otel"
	testingPackage = "testing"
	testifyPackage = "github.com/stretchr/testify/suite"
)

type basicGenerator struct {
	pkg           *model.Package
	src, dst      string
	filename      string
	interfaceName string
	withTests     bool

	packageMap map[string]string // map from import path to package name
}

func newGenerator(
	pkg *model.Package,
	src, dst string,
	interfaceName string,
	withTests bool,
) *basicGenerator {
	g := &basicGenerator{
		pkg:           pkg,
		src:           src,
		dst:           dst,
		filename:      "",
		interfaceName: interfaceName,
		withTests:     withTests,
	}

	// Get all required imports, and generate unique names for them all.
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

func (b *basicGenerator) generate() {
	for _, ifce := range b.pkg.Interfaces {
		if b.interfaceName == "" || b.interfaceName == ifce.Name {
			if b.withTests {
				b.generateImplementationTest(ifce)
			}
			b.generateImplementation(ifce)
		}
	}
}

func (b *basicGenerator) getArgNames(m *model.Method) []string {
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

func (b *basicGenerator) getArgTypes(m *model.Method, pkgOverride string) []string {
	argTypes := make([]string, len(m.In))
	for i, p := range m.In {
		argTypes[i] = p.Type.String(b.packageMap, pkgOverride)
	}
	if m.Variadic != nil {
		//nolint:makezero // mockgen authors
		argTypes = append(
			argTypes,
			"..."+m.Variadic.Type.String(b.packageMap, pkgOverride),
		)
	}
	return argTypes
}

func (b *basicGenerator) generateImplementation(ifce *model.Interface) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile(b.filename, protogen.GoImportPath(b.dst))
	path := b.dst + "/" + strtools.KebabCase(ifce.Name) + "/implementation.go"

	g.P("package ", strtools.SnakeCase(ifce.Name))
	g.P()
	g.P("type Implementation struct {")
	g.P("}")
	g.P()
	g.P("func NewImplementation() *", ifce.Name, "Implementation {")
	g.P("\treturn &", ifce.Name, "Implementation{}")
	g.P("}")

	err := fwriter.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}

	for _, method := range ifce.Methods {
		b.generateMethod(ifce, method)
	}
}

func (b *basicGenerator) generateImplementationTest(ifce *model.Interface) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile(b.filename, protogen.GoImportPath(b.dst))
	path := b.dst + "/" + strtools.KebabCase(ifce.Name) + "/implementation_test.go"

	g.P("package ", strtools.SnakeCase(ifce.Name))
	g.P()
	g.P("import (")
	for imp := range b.packageMap {
		g.P("\t\"", imp, "\"") // cannot use g.Import here
	}
	g.P("\t\"", testingPackage, "\"")
	g.P("\t\"", testifyPackage, "\"")
	g.P(")")
	g.P()
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

	err := fwriter.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}

	for _, method := range ifce.Methods {
		b.generateMethodTest(ifce, method)
	}
}

func (b *basicGenerator) generateMethod(ifce *model.Interface, m *model.Method) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile(b.filename, protogen.GoImportPath(b.dst))
	path := b.dst + "/" + strtools.KebabCase(ifce.Name) + "/" + strtools.SnakeCase(m.Name) + ".go"

	g.P("package ", strtools.SnakeCase(ifce.Name))
	g.P()

	g.P("import (")
	for imp := range b.packageMap {
		g.P("\t\"", imp, "\"") // cannot use g.Import here
	}
	g.P("\t\"", otelPackage, "\"")
	g.P(")")

	argNames := b.getArgNames(m)
	argTypes := b.getArgTypes(m, "")
	argString := mockgen.MakeArgString(argNames, argTypes)

	rets := make([]string, len(m.Out))
	for i, p := range m.Out {
		rets[i] = p.Type.String(b.packageMap, "")
	}
	retString := strings.Join(rets, ", ")
	if len(rets) > 1 {
		retString = "(" + retString + ")"
	}
	if retString != "" {
		retString = " " + retString
	}

	g.P(fmt.Sprintf("func (%v *%v%v) %v(%v)%v {", "i", ifce.Name, "Implementation", m.Name, argString,
		retString))

	if len(argTypes) > 0 && strings.HasPrefix(argTypes[0], "context.") {
		g.P(argNames[0], ", span := otel.Tracer(\"\").Start(", argNames[0], ", \"", ifce.Name, "Implementation.",
			m.Name, "\")")
		g.P("defer span.End()")
		g.P()
	}

	g.P("\tpanic(\"implement me\")")
	g.P("}")

	err := fwriter.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}
}

func (b *basicGenerator) generateMethodTest(ifce *model.Interface, m *model.Method) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile(b.filename, protogen.GoImportPath(b.dst))
	path := b.dst + "/" + strtools.KebabCase(ifce.Name) + "/" + strtools.SnakeCase(m.Name) + "_test.go"

	g.P("package ", strtools.SnakeCase(ifce.Name))
	g.P()

	g.P("import (")
	for imp := range b.packageMap {
		g.P("\t\"", imp, "\"") // cannot use g.Import here
	}
	g.P("\t\"", otelPackage, "\"")
	g.P(")")

	g.P("func (suite *TestSuite) Test", m.Name, "() {")
	g.P("}")

	err := fwriter.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}
}
