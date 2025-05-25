package in

import (
	"context"

	"github.com/not-for-prod/implgen/example/in/dto"
)

type TestInterface interface {
	A(ctx context.Context, req dto.GoRequest) error
	B(ctx context.Context, req map[dto.GoRequest]dto.GoRequest) error
	C(ctx context.Context, req []dto.GoRequest) error
	D(ctx context.Context, req int, opts ...dto.GoRequest) error
}
