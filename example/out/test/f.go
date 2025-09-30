package test

import (
	context "context"
	_ "embed"

	in "github.com/not-for-prod/implgen/example/in"
	"go.opentelemetry.io/otel"
)

func (i *Test) F(ctx context.Context, req in.FRequest) (in.FResponse, error) {
	ctx, span := otel.Tracer("my-brilliant-tracer").Start(ctx, "Test.F")
	defer span.End()

	panic("implement me")
}
