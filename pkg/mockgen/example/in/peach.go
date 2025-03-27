package in

import (
	"context"

	"github.com/google/uuid"
)

type Peach struct {
	Id    uuid.UUID
	Size  int
	Juice int
}

type PeachRepository interface {
	Create1(ctx context.Context, peach Peach) (uuid.UUID, error) // `sqlx:"GetContext"`
	Create2(peach Peach) (uuid.UUID, error)                      //sqlx:Get
	Create3(ctx context.Context, peach Peach) error              //sqlx:ExecContext
	Create4(peach Peach) error                                   //sqlx:Exec
}
