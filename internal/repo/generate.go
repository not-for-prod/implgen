package repo

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/not-for-prod/implgen/pkg/clog"
	"github.com/not-for-prod/implgen/pkg/fwriter"
	"github.com/not-for-prod/implgen/pkg/mockgen/model"
	"github.com/not-for-prod/implgen/pkg/strtools"
	"google.golang.org/protobuf/compiler/protogen"
)

type repoGenerator struct {
	src, dst, interfaceName string
	pkg                     *model.Package
	imports                 map[string]string // key: path, value: alias
	goPackage               string
}

func newRepoGenerator(src, dst, interfaceName string, pkg *model.Package) *repoGenerator {
	// modify imports with sql, otel && avito tx manager
	imports := map[string]string{
		"go.opentelemetry.io/otel":                          "otel",
		"github.com/jmoiron/sqlx":                           "sqlx",
		"github.com/avito-tech/go-transaction-manager/sqlx": "trmsqlx",
	}

	for imp := range pkg.Imports() {
		imports[imp] = filepath.Base(imp)
	}

	return &repoGenerator{
		src:           src,
		dst:           dst,
		interfaceName: interfaceName,
		pkg:           pkg,
		imports:       imports,
		goPackage:     getModuleName(),
	}
}

func (r *repoGenerator) packageName() string {
	return strings.ToLower(r.interfaceName) + "repo"
}

func (r *repoGenerator) generate() {
	for _, ifce := range r.pkg.Interfaces {
		if r.interfaceName == "" || r.interfaceName == ifce.Name {
			r.genInterface(ifce)
		}
	}
}

func (r *repoGenerator) genInterface(i *model.Interface) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile(r.dst, protogen.GoImportPath(r.dst))

	g.P("package ", r.packageName())
	g.P()
	g.P("import (")
	for path, alias := range r.imports {
		g.P("\t", alias, "\"", path, "\"") // cannot use g.Import here
	}
	g.P(")")
	g.P()
	g.P("type Implementation struct {")
	g.P("db *sqlx.DB")
	g.P("ctxGetter *trmsqlx.CtxGetter")
	g.P("}")
	g.P()
	g.P("func New(db *sqlx.DB, ctxGetter *trmsqlx.CtxGetter) *Implementation {")
	g.P("return &Implementation{")
	g.P("db: db,")
	g.P("ctxGetter: ctxGetter,")
	g.P("}")
	g.P("}")

	path := r.dst + "/" + strtools.KebabCase(i.Name) + "/implementation.go"

	err := fwriter.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}

	// generate sql embed file
	g = p.NewGeneratedFile(r.dst, protogen.GoImportPath(r.dst))
	g.P("package sql")
	g.P()
	g.P("import (")
	g.P("_ \"embed\"")
	g.P(")")
	g.P()

	for _, m := range i.Methods {
		g.P("//go:embed ", strtools.SnakeCase(m.Name), ".sql")
		g.P("var ", m.Name, " string")
		g.P()

		path = r.dst + "/" + strtools.KebabCase(i.Name) + "/sql/" + strtools.SnakeCase(m.Name) + ".sql"

		// create `.sql` file
		err = fwriter.WriteBytesToFile(path, []byte(``))
		if err != nil {
			clog.Fatal(err.Error())
		}

		r.genMethod(i, m)
	}

	path = r.dst + "/" + strtools.KebabCase(i.Name) + "/sql/sql.go"

	err = fwriter.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}
}

func (r *repoGenerator) genMethod(i *model.Interface, m *model.Method) {
	p := protogen.Plugin{}
	g := p.NewGeneratedFile(r.dst, protogen.GoImportPath(r.dst))

	g.P("package ", r.packageName())
	g.P()
	g.P("import (")
	for path, alias := range r.imports {
		g.P("\t", alias, "\"", path, "\"") // cannot use g.Import here
	}

	sqlPath := strings.ReplaceAll(r.dst, "./", "") + "/" + strtools.KebabCase(i.Name) + "/sql"

	g.P("sql \"", r.goPackage, "/", sqlPath, "\"")
	g.P(")")
	g.P()

	// get args and returns
	args := r.genArgs(m)
	returns := r.genReturns(m)
	returnString := r.genReturnString(m)

	g.P("func (i Implementation) ", m.Name, "(", args, ") (", returns, ") {")
	g.P("ctx, span := otel.Tracer(\"\").Start(ctx, \"", i.Name, "Implementation.", m.Name, "\")")
	g.P("defer span.End()")
	g.P()
	g.P("var err error")

	// parse comment to make propper sqlx method
	sqlxMethod := parseSQLXComment(m.Comment)

	switch sqlxMethod {
	case "ExecContext":
		g.P("err = i.ctxGetter.DefaultTrOrDB(ctx, i.db).ExecContext(ctx, sql.", m.Name, ")")
	case "SelectContext":
		g.P("var items []byte // TODO: fixit")
		g.P()
		g.P("err = i.ctxGetter.DefaultTrOrDB(ctx, i.db).SelectContext(ctx, &items, sql.", m.Name, ")")
	case "GetContext":
		g.P("var item []byte // TODO: fixit")
		g.P()
		g.P("err = i.ctxGetter.DefaultTrOrDB(ctx, i.db).GetContext(ctx, &item, sql.", m.Name, ")")
	default:
		g.P("err = i.ctxGetter.DefaultTrOrDB(ctx, i.db).ExecContext(ctx, )")
	}

	g.P("if err != nil {")
	g.P(returnString)
	g.P("}")
	g.P()

	g.P(returnString)
	g.P("}")
	g.P()

	path := r.dst + "/" + strtools.KebabCase(i.Name) + "/" + strtools.SnakeCase(m.Name) + ".go"

	err := fwriter.WriteGeneratedFile(path, g)
	if err != nil {
		clog.Fatal(err.Error())
	}
}

func (r *repoGenerator) genReturns(m *model.Method) string {
	builder := strings.Builder{}

	for _, p := range m.Out {
		builder.WriteString(p.Type.String(r.imports, "") + ", ")
	}

	return builder.String()
}

// genReturnString is vibecoded)))
func (r *repoGenerator) genReturnString(m *model.Method) string {
	builder := strings.Builder{}

	for i, p := range m.Out {
		if i > 0 {
			builder.WriteString(", ")
		}

		returnType := p.Type.String(r.imports, "")

		// Handle different return types
		switch returnType {
		case "error":
			builder.WriteString("err")
		case "string", "model.OrderID": // If it's a string-based type
			builder.WriteString(`""`)
		case "int", "int32", "int64":
			builder.WriteString("0")
		case "bool":
			builder.WriteString("false")
		case "float32", "float64":
			builder.WriteString("0.0")
		default:
			// here we need to add it in model
			p := protogen.Plugin{}
			g := p.NewGeneratedFile(r.dst, protogen.GoImportPath(r.dst))
			g.P("package model")
			g.P()
			tokens := strings.Split(returnType, ".")
			modelName := tokens[len(tokens)-1]
			g.P("type ", modelName, " struct {}")

			path := r.dst + "/" + strtools.KebabCase(r.interfaceName) + "/model/" +
				strtools.SnakeCase(modelName) + ".go"

			err := fwriter.WriteGeneratedFile(path, g)
			if err != nil {
				clog.Fatal(err.Error())
			}

			if strings.HasPrefix(returnType, "*") || strings.HasPrefix(returnType, "[]") {
				builder.WriteString("nil") // Pointers and slices should return nil
			} else {
				builder.WriteString(returnType + "{}") // Structs return empty struct
			}
		}
	}

	return "return " + builder.String()
}

func (r *repoGenerator) genArgs(m *model.Method) string {
	builder := strings.Builder{}

	for i, p := range m.In {
		name := p.Name
		if name == "" || name == "_" {
			name = fmt.Sprintf("arg%d", i)
		}

		_type := p.Type.String(r.imports, "")

		builder.WriteString(name + " " + _type + ",")
	}
	if m.Variadic != nil {
		name := m.Variadic.Name
		if name == "" {
			name = fmt.Sprintf("arg%d", len(m.In))
		}

		builder.WriteString(name + " " + m.Variadic.Type.String(r.imports, "") + ",")
	}
	return builder.String()
}
