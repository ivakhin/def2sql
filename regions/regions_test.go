package regions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadRegions(t *testing.T) {
	res, err := Load("./testdata/null.yml")
	if !assert.Error(t, err) {
		return
	}

	res, err = Load("./testdata/regions_test.yml")
	if !assert.Nilf(t, err, "%+v", err) {
		return
	}

	expected := Regions{
		1: {"A", []string{"aaa"}, nil},
		2: {"B", []string{"bbb"}, nil},
		3: {"C", []string{"ccc"}, nil},
	}

	if !assert.Equal(t, expected, res) {
		return
	}
}

func TestRegions_Match(t *testing.T) {
	r := Regions{
		1: {"A", []string{"aaa"}, []string{"A"}},
		2: {"B", []string{"bbb", "c"}, nil},
		3: {"C", []string{"ccc"}, nil},
	}

	data := []struct {
		actual   string
		expected map[int]struct{}
	}{
		{
			"aaa",
			map[int]struct{}{1: {}},
		},
		{
			"ccc",
			map[int]struct{}{2: {}, 3: {}},
		},
		{
			"aaaAccc",
			map[int]struct{}{2: {}, 3: {}},
		},
	}

	for _, d := range data {
		if !assert.Equal(t, d.expected, r.Match(d.actual)) {
			return
		}
	}
}
