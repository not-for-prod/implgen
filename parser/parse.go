package parser

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/not-for-prod/implgen/model"
	"golang.org/x/tools/go/packages"
)

type Command struct {
	src        string
	info       *types.Info
	selfImport model.Import
	imports    []model.Import
}

func NewCommand(src string) *Command {
	return &Command{
		src: src,
	}
}

func (cmd *Command) Execute() (model.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedImports,
		Dir: filepath.Dir(cmd.src), // important: evaluate from your file’s directory
	}

	abs, err := filepath.Abs(cmd.src)
	if err != nil {
		return model.Package{}, err
	}

	pkgs, err := packages.Load(cfg, fmt.Sprintf("file=%s", abs))
	if err != nil {
		return model.Package{}, fmt.Errorf("failed to load package: %w", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return model.Package{}, fmt.Errorf("package load errors")
	}

	var astFile *ast.File
	var pkg *packages.Package

	for _, p := range pkgs {
		for _, f := range p.Syntax {
			// filename matching
			if p.Fset.File(f.Pos()).Name() == abs {
				astFile = f
				pkg = p
				break
			}
		}
	}

	if pkg == nil {
		return model.Package{}, fmt.Errorf("package not found")
	}

	// add self import
	selfImport := model.Import{
		Alias: astFile.Name.Name, // package name as alias
		Path:  packageImportPath(),
	}
	cmd.imports = append(cmd.imports, selfImport)

	cmd.info = pkg.TypesInfo
	cmd.imports = append(cmd.imports, cmd.parseImports(astFile)...)

	return model.Package{
		Name:       pkg.Name,
		Interfaces: cmd.parseInterfaces(astFile),
		Imports:    cmd.imports,
	}, nil
}

func (cmd *Command) parseImports(node *ast.File) []model.Import {
	var imports []model.Import

	for _, imp := range node.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			alias = filepath.Base(path)
		}

		imports = append(
			imports, model.Import{
				Alias: alias,
				Path:  path,
			},
		)
	}

	return imports
}

func (cmd *Command) parseInterfaces(node *ast.File) []model.Interface {
	var interfaces []model.Interface

	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			tspec := spec.(*ast.TypeSpec)
			if iface, ok := tspec.Type.(*ast.InterfaceType); ok {
				interfaces = append(interfaces, cmd.parseInterface(tspec.Name.Name, iface))
			}
		}
	}
	return interfaces
}

func (cmd *Command) parseInterface(name string, iface *ast.InterfaceType) model.Interface {
	var methods []model.Method

	for _, field := range iface.Methods.List {
		if len(field.Names) == 0 {
			continue // skip embedded interfaces
		}
		methodName := field.Names[0].Name
		ftype, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		methods = append(methods, cmd.parseMethod(methodName, ftype))
	}

	return model.Interface{
		Name:    name,
		Methods: methods,
	}
}

func (cmd *Command) parseMethod(name string, ftype *ast.FuncType) model.Method {
	method := model.Method{Name: name}

	if ftype.Params != nil {
		for i, param := range ftype.Params.List {
			typ := cmd.exprString(param.Type)
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
			typ := cmd.exprString(result.Type)
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

func (cmd *Command) exprString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		// look up type info
		if obj, ok := cmd.info.Uses[e]; ok {
			switch obj := obj.(type) {
			case *types.TypeName:
				pkg := obj.Pkg()
				if pkg == nil {
					// builtin like int, string, error
					return e.Name
				}
				if pkg.Name() == cmd.selfImport.Alias {
					// local type
					return cmd.selfImport.Alias + "." + e.Name
				}
				return pkg.Name() + "." + e.Name
			default:
				return e.Name
			}
		}
		return e.Name
	case *ast.SelectorExpr:
		return cmd.exprString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + cmd.exprString(e.X)
	case *ast.ArrayType:
		return "[]" + cmd.exprString(e.Elt)
	case *ast.Ellipsis:
		return "..." + cmd.exprString(e.Elt)
	case *ast.MapType:
		return "map[" + cmd.exprString(e.Key) + "]" + cmd.exprString(e.Value)
	case *ast.FuncType:
		return "func" // Simplified
	default:
		return "unknown"
	}
}

func packageImportPath() string {
	// Get path where implgen is executed
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	modRoot, _ := findGoModRoot(pwd)
	if modRoot != "" {
		// go.mod module name
		modName, _ := moduleName(modRoot)
		dir, _ := filepath.Rel(modRoot, pwd)

		return modName + "/" + dir
	}

	// Option B: fallback to directory name
	return filepath.Base(pwd)
}

func moduleName(modRoot string) (string, error) {
	f, err := os.Open(filepath.Join(modRoot, "go.mod"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", nil
}

// findGoModRoot get go mod absolute path
func findGoModRoot(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		modPath := filepath.Join(absDir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			return absDir, nil
		}

		parent := filepath.Dir(absDir)
		if parent == absDir { // reached filesystem root
			break
		}
		absDir = parent
	}

	return "", nil
}
