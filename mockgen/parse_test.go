package mockgen

import (
	"testing"

	"github.com/not-for-prod/implgen/pkg/clog"
)

// TestSourceMode was supposed for testing in debug mode
func TestSourceMode(t *testing.T) {
	source := "./example/in/peach.go"

	_package, err := SourceMode(source)
	if err != nil {
		clog.Error(err.Error())
	}

	clog.Info(_package)
}
