package filewriter

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// WriteToFile writes r to the file with path
func WriteToFile(path string, r io.Reader) error {
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
