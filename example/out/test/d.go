package test

import (
	context "context"
	_ "embed"

	dto "github.com/not-for-prod/implgen/example/in/dto"
)

func (i *Test) D(ctx context.Context, req int, opts ...dto.GoRequest) error {
	panic("implement me")
}
