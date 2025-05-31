package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/not-for-prod/implgen/model"
)

type ParseCommand struct {
	src string
}

func New(src string) *ParseCommand {
	return &ParseCommand{
		src: src,
	}
}

func (cmd *ParseCommand) Execute() (model.Package, error) {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, cmd.src, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return model.Package{}, fmt.Errorf("failed parsing source file %v: %v", cmd.src, err)
	}

	pkg := model.Package{
		Name:       file.Name.Name,
		Interfaces: parseInterfaces(file),
		Imports:    parseImports(file),
	}

	return pkg, nil
}

func parseImports(node *ast.File) []model.Import {
	imports := make([]model.Import, 0)
	for _, imp := range node.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			alias = filepath.Base(path)
		}
		imports = append(imports, model.Import{
			Alias: alias,
			Path:  path,
		})
	}

	return imports
}

func parseInterfaces(node *ast.File) []model.Interface {
	var interfaces []model.Interface

	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			tspec := spec.(*ast.TypeSpec)
			if iface, ok := tspec.Type.(*ast.InterfaceType); ok {
				interfaces = append(interfaces, parseInterface(tspec.Name.Name, iface))
			}
		}
	}
	return interfaces
}

func parseInterface(name string, iface *ast.InterfaceType) model.Interface {
	methods := []model.Method{}

	for _, field := range iface.Methods.List {
		if len(field.Names) == 0 {
			continue // skip embedded interfaces
		}
		methodName := field.Names[0].Name
		ftype, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		methods = append(methods, parseMethod(methodName, ftype))
	}

	return model.Interface{
		Name:    name,
		Methods: methods,
	}
}

func parseMethod(name string, ftype *ast.FuncType) model.Method {
	method := model.Method{Name: name}

	if ftype.Params != nil {
		for i, param := range ftype.Params.List {
			typ := exprString(param.Type)
			if len(param.Names) == 0 {
				method.In = append(method.In, model.Parameter{Name: "arg" + string(rune(i+'a')), Type: typ})
			} else {
				for _, name := range param.Names {
					method.In = append(method.In, model.Parameter{Name: name.Name, Type: typ})
				}
			}
		}
	}

	if ftype.Results != nil {
		for i, result := range ftype.Results.List {
			typ := exprString(result.Type)
			if len(result.Names) == 0 {
				method.Out = append(method.Out, model.Parameter{Name: "ret" + string(rune(i+'a')), Type: typ})
			} else {
				for _, name := range result.Names {
					method.Out = append(method.Out, model.Parameter{Name: name.Name, Type: typ})
				}
			}
		}
	}

	return method
}

func exprString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + exprString(e.X)
	case *ast.ArrayType:
		return "[]" + exprString(e.Elt)
	case *ast.Ellipsis:
		return "..." + exprString(e.Elt)
	case *ast.MapType:
		return "map[" + exprString(e.Key) + "]" + exprString(e.Value)
	case *ast.FuncType:
		return "func" // Simplified
	default:
		return "unknown"
	}
}
