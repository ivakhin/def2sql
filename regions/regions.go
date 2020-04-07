package regions

import (
	"strings"

	"def2sql/helpers"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Regions map[int]Region

type Region struct {
	Name       string   `yaml:"name"`
	Contain    []string `yaml:"contain"`
	NotContain []string `yaml:"not_contain"`
}

func Load(path string) (Regions, error) {
	data, err := helpers.Read(path)
	if err != nil {
		return nil, err
	}

	regs := Regions{}
	err = yaml.Unmarshal(data, &regs)
	return regs, errors.WithStack(err)
}

func (r Regions) Match(source string) map[int]struct{} {
	overlaps := make(map[int]struct{}, 0)
	for code, reg := range r {
	loop:
		for _, contain := range reg.Contain {
			if strings.Contains(source, contain) {
				for _, notContain := range reg.NotContain {
					if strings.Contains(source, notContain) {
						break loop
					}
				}
				overlaps[code] = struct{}{}
			}
		}
	}

	return overlaps
}
