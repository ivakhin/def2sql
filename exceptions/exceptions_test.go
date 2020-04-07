package exceptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadExceptions(t *testing.T) {
	res, err := Load("./testdata/null.yml")
	if !assert.Error(t, err) {
		return
	}

	res, err = Load("./testdata/exceptions_test.yml")
	if !assert.Nilf(t, err, "%+v", err) {
		return
	}

	expected := Exceptions{
		"A": "B",
		"C": "D",
	}

	if !assert.Equal(t, expected, res) {
		return
	}
}

func TestExceptions_Apply(t *testing.T) {
	e := Exceptions{
		"A": "a",
	}
	data := []byte("Abc")
	expected := []byte("abc")

	if !assert.Equal(t, expected, e.Apply(data)) {
		return
	}
}
