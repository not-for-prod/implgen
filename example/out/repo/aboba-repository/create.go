package abobarepositoryrepo

import (
	context "context"

	in "github.com/not-for-prod/implgen/example/in"
	model "github.com/not-for-prod/implgen/example/in/model"
	sql "github.com/not-for-prod/implgen/example/out/repo/aboba-repository/sql"
	otel "go.opentelemetry.io/otel"
)

func (i Implementation) Create(ctx context.Context, req in.CreateRequest) (model.OrderID, error) {
	ctx, span := otel.Tracer("").Start(ctx, "AbobaRepositoryImplementation.Create")
	defer span.End()

	var err error
	var item []byte // TODO: fixit

	err = i.ctxGetter.DefaultTrOrDB(ctx, i.db).GetContext(ctx, &item, sql.Create)
	if err != nil {
		return "", err
	}

	return "", err
}
