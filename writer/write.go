package writer

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/not-for-prod/implgen/model"
	"github.com/not-for-prod/implgen/pkg/clog"

	importsTool "golang.org/x/tools/imports"
)

type Command struct {
	// Enable verbose logging
	verbose bool
}

func NewCommand(verbose bool) *Command {
	return &Command{verbose: verbose}
}

func (w *Command) overwrite(path string) bool {
	if !w.verbose {
		return false
	}

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

// WriteToFile writes r to the file with path
func (w *Command) writeToFile(path string, r io.Reader) error {
	if !w.overwrite(path) {
		return nil
	}

	dir := filepath.Dir(path)
	if dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, r)
	return err
}

func (w *Command) writeBytesToFile(path string, data []byte) error {
	return w.writeToFile(path, bytes.NewReader(data))
}

// writeGoBytesToFile WriteBytesToFile (that are actually go code)  but before makes `goimports -w ...` && `go fmt ...`
func (w *Command) writeGoBytesToFile(path string, data []byte) error {
	var err error

	// goimports -w ...
	data, err = importsTool.Process(
		path, data, &importsTool.Options{
			Comments: true,
		},
	)
	if err != nil {
		clog.Fatalf(err.Error())
	}

	// go fmt ...
	data, err = format.Source(data)
	if err != nil {
		clog.Fatalf(err.Error())
	}

	err = w.writeBytesToFile(path, data)
	if err != nil {
		return err
	}

	return nil
}

func (w *Command) Execute(files []model.File) error {
	for _, file := range files {
		err := w.writeGoBytesToFile(file.Path, file.Data)
		
		return err
	}

	return nil
}
