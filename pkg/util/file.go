package util

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}
