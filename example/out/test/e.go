package test

import (
	context "context"
	_ "embed"

	in "github.com/not-for-prod/implgen/example/in"
	"go.opentelemetry.io/otel"
)

func (i *Test) E(ctx context.Context, req in.ERequest) (in.EResponse, error) {
	ctx, span := otel.Tracer("my-brilliant-tracer").Start(ctx, "Test.E")
	defer span.End()

	panic("implement me")
}
