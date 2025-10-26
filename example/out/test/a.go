package test

import (
	context "context"
	_ "embed"

	dto "github.com/not-for-prod/implgen/example/in/dto"
)

func (i *Test) A(ctx context.Context, req dto.GoRequest) error {
	panic("implement me")
}
