package test_interface

import (
	context "context"
	_ "embed"

	in "github.com/not-for-prod/implgen/example/in"
)

func (i *Implementation) E(ctx context.Context, req in.ERequest) (in.EResponse, error) {
	panic("implement me")
}
