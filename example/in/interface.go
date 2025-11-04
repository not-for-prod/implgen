package in

import (
	"context"
	_ "embed"

	"github.com/not-for-prod/implgen/example/in/dto"
)

//go:generate implgen --src interface.go --interface-name TestInterface --dst ../test
//go:generate pwd

type ERequest struct{}

type EResponse struct{}

type TestInterface interface {
	A(ctx context.Context, req dto.GoRequest) error
	B(ctx context.Context, req map[dto.GoRequest]dto.GoRequest) error
	C(ctx context.Context, req []dto.GoRequest) error
	D(ctx context.Context, req int, opts ...dto.GoRequest) error
	E(ctx context.Context, req ERequest) (EResponse, error)
	F(ctx context.Context, req FRequest) (FResponse, error)
}
