package test

import (
	context "context"

	dto "github.com/not-for-prod/implgen/example/in/dto"
	"go.opentelemetry.io/otel"
)

func (i *Test) C(ctx context.Context, req []dto.GoRequest) error {
	ctx, span := otel.Tracer("my-brilliant-tracer").Start(ctx, "Test.C")
	defer span.End()

	panic("implement me")
}
