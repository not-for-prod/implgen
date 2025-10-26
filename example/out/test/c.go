package test

import (
	context "context"
	_ "embed"

	dto "github.com/not-for-prod/implgen/example/in/dto"
)

func (i *Test) C(ctx context.Context, req []dto.GoRequest) error {
	panic("implement me")
}
