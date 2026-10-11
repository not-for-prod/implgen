package test

import (
	context "context"
	_ "embed"

	in "github.com/not-for-prod/implgen/example/in"
)

func (i *Test) F(ctx context.Context, req in.FRequest) (in.FResponse, error) {
	panic("implement me")
}
