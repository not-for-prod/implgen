package test

import (
	context "context"
	_ "embed"

	dto "github.com/not-for-prod/implgen/example/in/dto"
)

func (i *Test) B(ctx context.Context, req map[dto.GoRequest]dto.GoRequest) error {
	panic("implement me")
}
