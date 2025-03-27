package aboba_repositoryrepo

import (
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	sqlx "github.com/jmoiron/sqlx"
)

type Implementation struct {
	db        *sqlx.DB
	ctxGetter *trmsqlx.CtxGetter
}

func New(db *sqlx.DB, ctxGetter *trmsqlx.CtxGetter) *Implementation {
	return &Implementation{
		db:        db,
		ctxGetter: ctxGetter,
	}
}
