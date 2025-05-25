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
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/not-for-prod/implgen/pkg/clog"

	importsTool "golang.org/x/tools/imports"
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

// WriteToFile writes r to the file with path
func WriteToFile(path string, r io.Reader) error {
	if !overwrite(path) {
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

func WriteStringToFile(path string, s string) error {
	return WriteToFile(path, strings.NewReader(s))
}

func WriteBytesToFile(path string, data []byte) error {
	return WriteToFile(path, bytes.NewReader(data))
}

// WriteGoBytesToFile WriteBytesToFile (that are actually go code)  but before makes `goimports -w ...` && `go fmt ...`
func WriteGoBytesToFile(path string, data []byte) error {
	var err error

	// goimports -w ...
	data, err = importsTool.Process(path, data, &importsTool.Options{
		Comments: true,
	})
	if err != nil {
		clog.Fatalf(err.Error())
	}

	//// go fmt ...
	data, err = format.Source(data)
	if err != nil {
		clog.Fatalf(err.Error())
	}

	err = WriteBytesToFile(path, data)
	if err != nil {
		return err
	}

	return nil
}
