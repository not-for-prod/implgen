package implgen

import (
	"go/format"

	"github.com/not-for-prod/gen-tools/clog"
	"github.com/not-for-prod/gen-tools/fwriter"
	importsTool "golang.org/x/tools/imports"
	"google.golang.org/protobuf/compiler/protogen"
)

func writeGeneratedFile(g *protogen.GeneratedFile, filepath string) {
	// Proceed with writing the file
	data, err := g.Content()
	if err != nil {
		clog.Fatalf(err.Error())
	}

	// goimports -w ...
	data, err = importsTool.Process(filepath, data, &importsTool.Options{})
	if err != nil {
		clog.Fatalf(err.Error())
	}

	// go fmt ...
	data, err = format.Source(data)
	if err != nil {
		clog.Fatalf(err.Error())
	}

	err = fwriter.WriteStringToFile(filepath, string(data))
	if err != nil {
		clog.Fatalf(err.Error())
	}

	clog.Info("file", filepath, "written")
}
