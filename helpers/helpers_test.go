package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isHTTPSource(t *testing.T) {
	src := "https://example.com"
	if !assert.True(t, isHTTPSource(src)) {
		return
	}

	src = "./file"
	if !assert.False(t, isHTTPSource(src)) {
		return
	}
}
