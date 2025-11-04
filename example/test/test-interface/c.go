package test_interface

import (
	context "context"
	_ "embed"

	dto "github.com/not-for-prod/implgen/example/in/dto"
)

func (i *Implementation) C(ctx context.Context, req []dto.GoRequest) error {
	panic("implement me")
}
