package config

import (
	"def2sql/helpers"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DataSource []string `yaml:"data_source"`
	Exceptions string   `yaml:"exceptions"`
	Regions    string   `yaml:"regions"`
	DB         db       `yaml:"db"`
}

type db struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	Table    string `yaml:"table"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func Read(path string) (*Config, error) {
	data, err := helpers.Read(path)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = yaml.Unmarshal(data, &c)
	return &c, errors.WithStack(err)
}
