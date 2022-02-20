package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFilterFile(t *testing.T) {
	assert := assert.New(t)

	var rf RegexpFilter
	rf.exprList = make([]*regexp.Regexp, 0)

	rf.readFilterFile("./tests/blacklist.1")
	assert.Equal(len(rf.exprList), 11)
	assert.True(rf.IsMatch("adtracking.foo.com"))
	assert.False(rf.IsMatch("foo.com"))

	rf.readFilterFile("./tests/blacklist.2")
	assert.Equal(len(rf.exprList), 13)
	assert.True(rf.IsMatch("www.yandex.ru"))
	assert.False(rf.IsMatch("www.yandex.com"))
}

func TestIsFiltered(t *testing.T) {
	assert := assert.New(t)

	var fd FilteredDomains
	fd.init()

	fd.whiteList.readFilterFile("./tests/whitelist.1")
	fd.blackList.readFilterFile("./tests/blacklist.1")
	fd.blackList.readFilterFile("./tests/blacklist.2")

	assert.True(fd.isFiltered("adtracking.foo.com"))
	assert.False(fd.isFiltered("foo.com"))
	assert.False(fd.isFiltered("www.yandex.ru"))
	assert.True(fd.isFiltered("www.foo.ru"))
}
