package implgen

import (
	"fmt"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/not-for-prod/implgen/pkg/mockgen/model"
	"github.com/not-for-prod/implgen/pkg/string-tools"
	"google.golang.org/protobuf/compiler/protogen"
)

const otelPackage = "go.opentelemetry.io/otel"

type Generator struct {
	pkg           *model.Package
	src, dst      string
	filename      string
	withOtel      bool
	interfaceName string

	packageMap map[string]string // map from import path to package name
}

func NewGenerator(
	pkg *model.Package,
	src, dst string,
	withOtel bool,
	interfaceName string,
) *Generator {
	g := &Generator{
		pkg:           pkg,
		src:           src,
		dst:           dst,
		filename:      "",
		withOtel:      withOtel,
		interfaceName: interfaceName,
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

	packagesName := createPackageMap(sortedPaths)

	g.packageMap = make(map[string]string, len(im))
	localNames := make(map[string]bool, len(im))
	for _, pth := range sortedPaths {
		base, ok := packagesName[pth]
		if !ok {
			base = sanitize(path.Base(pth))
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

func (gen *Generator) Generate() {
	for _, ifce := range gen.pkg.Interfaces {
		if gen.interfaceName == "" || gen.interfaceName == ifce.Name {
			gen.generateInterface(ifce)
		}
	}
}

func (gen *Generator) getArgNames(m *model.Method) []string {
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

func (gen *Generator) getArgTypes(m *model.Method, pkgOverride string) []string {
	argTypes := make([]string, len(m.In))
	for i, p := range m.In {
		argTypes[i] = p.Type.String(gen.packageMap, pkgOverride)
	}
	if m.Variadic != nil {
		//nolint:makezero // mockgen authors
		argTypes = append(
			argTypes,
			"..."+m.Variadic.Type.String(gen.packageMap, pkgOverride),
		)
	}
	return argTypes
}

func (gen *Generator) generateInterface(ifce *model.Interface) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile(gen.filename, protogen.GoImportPath(gen.dst))

	g.P("package ", strtools.SnakeCase(gen.pkg.Name))
	g.P()
	g.P("type ", ifce.Name, "Implementation struct {")
	g.P("}")
	g.P()
	g.P("func NewGenerator", ifce.Name, "Implementation() *", ifce.Name, "Implementation {")
	g.P("\treturn &", ifce.Name, "Implementation{}")
	g.P("}")

	writeGeneratedFile(g, gen.dst+"/"+strtools.KebabCase(ifce.Name)+"/implementation.go")

	for _, method := range ifce.Methods {
		gen.generateMethod(ifce, method)
	}
}

func (gen *Generator) generateMethod(ifce *model.Interface, m *model.Method) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile(gen.filename, protogen.GoImportPath(gen.dst))

	g.P("package ", strtools.SnakeCase(gen.pkg.Name))
	g.P()

	g.P("import (")
	for imp := range gen.packageMap {
		g.P("\t\"", imp, "\"") // cannot use g.Import here
	}
	if gen.withOtel {
		g.P("\t\"", otelPackage, "\"")
	}
	g.P(")")

	argNames := gen.getArgNames(m)
	argTypes := gen.getArgTypes(m, "")
	argString := makeArgString(argNames, argTypes)

	rets := make([]string, len(m.Out))
	for i, p := range m.Out {
		rets[i] = p.Type.String(gen.packageMap, "")
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

	if len(argTypes) > 0 && strings.HasPrefix(argTypes[0], "context.") && gen.withOtel {
		g.P("ctx, span := otel.Tracer(\"\").Start(", argNames[0], ", \"", ifce.Name, "Implementation.", m.Name, "\")")
		g.P("defer span.End()")
		g.P()
	}

	g.P("\tpanic(\"implement me\")")
	g.P("}")

	writeGeneratedFile(g, gen.dst+"/"+strtools.KebabCase(ifce.Name)+"/"+strtools.SnakeCase(m.Name)+".go")
}
