package exceptions

import (
	"bytes"

	"def2sql/helpers"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Exceptions map[string]string

func Load(path string) (Exceptions, error) {
	data, err := helpers.Read(path)
	if err != nil {
		return nil, err
	}

	ex := Exceptions{}
	err = yaml.Unmarshal(data, &ex)
	return ex, errors.WithStack(err)
}

func (e Exceptions) Apply(data []byte) []byte {
	for a, b := range e {
		data = bytes.ReplaceAll(data, []byte(a), []byte(b))
	}
	return data
}
