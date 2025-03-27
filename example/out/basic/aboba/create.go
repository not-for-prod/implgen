package in

import (
	"context"

	"github.com/not-for-prod/implgen/example/in"
	"github.com/not-for-prod/implgen/example/in/model"
	"go.opentelemetry.io/otel"
)

func (i *AbobaImplementation) Create(ctx context.Context, req in.CreateRequest) (model.OrderID, error) {
	ctx, span := otel.Tracer("").Start(ctx, "AbobaImplementation.Create")
	defer span.End()

	panic("implement me")
}
