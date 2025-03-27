package in

import (
	"context"

	"github.com/not-for-prod/implgen/example/in/model"
	"go.opentelemetry.io/otel"
)

func (i *AbobaImplementation) Get(ctx context.Context, id model.OrderID) (model.Order, error) {
	ctx, span := otel.Tracer("").Start(ctx, "AbobaImplementation.Get")
	defer span.End()

	panic("implement me")
}
