package test_interface

import (
	context "context"
	_ "embed"

	in "github.com/not-for-prod/implgen/example/in"
)

func (i *Implementation) F(ctx context.Context, req in.FRequest) (in.FResponse, error) {
	panic("implement me")
}
