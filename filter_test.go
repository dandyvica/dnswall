package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFilterFile(t *testing.T) {
	assert := assert.New(t)

	blocklist := "./tests/blocklist.1"
	blockRegex := ReadFilterFile(blocklist, 10)

	assert.Equal(blockRegex.filterFile, blocklist)
	assert.Equal(len(blockRegex.exprList), 11)

	assert.True(blockRegex.IsMatch("adtracking.foo.com"))
	assert.False(blockRegex.IsMatch("foo.com"))

}
