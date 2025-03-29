package abobarepositoryrepo

import (
	context "context"

	model "github.com/not-for-prod/implgen/example/in/model"
	sql "github.com/not-for-prod/implgen/example/out/repo/aboba-repository/sql"
	otel "go.opentelemetry.io/otel"
)

func (i Implementation) Get(ctx context.Context, id model.OrderID) (model.Order, error) {
	ctx, span := otel.Tracer("").Start(ctx, "AbobaRepositoryImplementation.Get")
	defer span.End()

	var err error
	var item []byte // TODO: fixit

	err = i.ctxGetter.DefaultTrOrDB(ctx, i.db).GetContext(ctx, &item, sql.Get)
	if err != nil {
		return model.Order{}, err
	}

	return model.Order{}, err
}
