package service

import (
	"testing"

	"def2sql/config"

	"github.com/stretchr/testify/assert"
)

func TestService_get(t *testing.T) {
	c, err := config.Read("testdata/config.yml")
	if !assert.Nil(t, err) {
		return
	}

	s, err := New(c)
	if !assert.Nil(t, err) {
		return
	}

	data, err := s.get()
	if !assert.Nil(t, err) {
		return
	}

	rec := s.read(data)

	if !assert.Equal(t, 0, len(rec.Wrong)) {
		return
	}

	if !assert.Equal(t, 0, len(rec.UnknownRegion)) {
		return
	}

	if !assert.Equal(t, 0, len(rec.Warnings)) {
		return
	}
}
