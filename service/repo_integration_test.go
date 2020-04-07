// +build integration

package service

import (
	"testing"

	"def2sql/config"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

func Test_repo_truncate(t *testing.T) {
	c, err := config.Read("testdata/config.yml")
	if !assert.Nil(t, err) {
		return
	}

	s, err := New(c)
	if !assert.Nil(t, err) {
		return
	}

	if !assert.Nil(t, s.repo.truncate()) {
		return
	}
}

func Test_repo_create(t *testing.T) {
	c, err := config.Read("testdata/config.yml")
	if !assert.Nil(t, err) {
		return
	}

	s, err := New(c)
	if !assert.Nil(t, err) {
		return
	}

	if !assert.Nil(t, s.repo.create()) {
		return
	}
}
