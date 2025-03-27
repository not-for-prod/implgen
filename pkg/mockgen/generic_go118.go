package mockgen

import (
	"go/ast"

	"github.com/not-for-prod/implgen/pkg/mockgen/model"
)

func getTypeSpecTypeParams(ts *ast.TypeSpec) []*ast.Field {
	if ts == nil || ts.TypeParams == nil {
		return nil
	}
	return ts.TypeParams.List
}

func (p *fileParser) parseGenericType(pkg string, typ ast.Expr, tps map[string]bool) (model.Type, error) {
	switch v := typ.(type) {
	case *ast.IndexExpr:
		m, err := p.parseType(pkg, v.X, tps)
		if err != nil {
			return nil, err
		}
		nm, ok := m.(*model.NamedType)
		if !ok {
			return m, nil
		}
		t, err := p.parseType(pkg, v.Index, tps)
		if err != nil {
			return nil, err
		}
		nm.TypeParams = &model.TypeParametersType{TypeParameters: []model.Type{t}}
		return m, nil
	case *ast.IndexListExpr:
		m, err := p.parseType(pkg, v.X, tps)
		if err != nil {
			return nil, err
		}
		nm, ok := m.(*model.NamedType)
		if !ok {
			return m, nil
		}
		var ts []model.Type
		for _, expr := range v.Indices {
			t, err := p.parseType(pkg, expr, tps)
			if err != nil {
				return nil, err
			}
			ts = append(ts, t)
		}
		nm.TypeParams = &model.TypeParametersType{TypeParameters: ts}
		return m, nil
	}
	return nil, nil
}
