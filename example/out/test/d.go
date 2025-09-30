package test

import (
	context "context"
	_ "embed"

	dto "github.com/not-for-prod/implgen/example/in/dto"
	"go.opentelemetry.io/otel"
)

func (i *Test) D(ctx context.Context, req int, opts ...dto.GoRequest) error {
	ctx, span := otel.Tracer("my-brilliant-tracer").Start(ctx, "Test.D")
	defer span.End()

	panic("implement me")
}
