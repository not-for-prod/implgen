package implgen

import (
	"errors"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	file_writer "github.com/not-for-prod/implgen/pkg/file-writer"
	"github.com/not-for-prod/implgen/pkg/logger"
	importsTool "golang.org/x/tools/imports"
	"google.golang.org/protobuf/compiler/protogen"
)

func overwrite(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return true
	}

	keep := "don't overwrite " + path
	write := "overwrite " + path

	prompt := promptui.Select{
		Label: fmt.Sprintf("`%s` already exists, do you want to keep this file", filepath.Base(path)),
		Items: []string{keep, write},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false
	}

	if result == write {
		return true
	}

	return false
}

func writeGeneratedFile(g *protogen.GeneratedFile, filepath string) {
	// Check if the file already exists
	if !overwrite(filepath) {
		return
	}

	// Proceed with writing the file
	data, err := g.Content()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	// goimports -w ...
	data, err = importsTool.Process(filepath, data, &importsTool.Options{})
	if err != nil {
		logger.Fatalf(err.Error())
	}

	// go fmt ...
	data, err = format.Source(data)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	err = file_writer.WriteStringToFile(filepath, string(data))
	if err != nil {
		logger.Fatalf(err.Error())
	}

	logger.Info("file", filepath, "written")
}
