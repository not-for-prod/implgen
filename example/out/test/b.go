package test

import (
	context "context"
	_ "embed"

	dto "github.com/not-for-prod/implgen/example/in/dto"
	"go.opentelemetry.io/otel"
)

func (i *Test) B(ctx context.Context, req map[dto.GoRequest]dto.GoRequest) error {
	ctx, span := otel.Tracer("my-brilliant-tracer").Start(ctx, "Test.B")
	defer span.End()

	panic("implement me")
}
