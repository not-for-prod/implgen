package test_interface

import (
	context "context"
	_ "embed"

	dto "github.com/not-for-prod/implgen/example/in/dto"
)

func (i *Implementation) B(ctx context.Context, req map[dto.GoRequest]dto.GoRequest) error {
	panic("implement me")
}
