package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpack(t *testing.T) {
	res, err := unpack("")
	assert.NoError(t, err)
	assert.Equal(t, "", res)

	res, err = unpack("abcd")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", res)

	res, err = unpack("a4bc2d5e")
	assert.NoError(t, err)
	assert.Equal(t, "aaaabccddddde", res)

	res, err = unpack("a4c2")
	assert.NoError(t, err)
	assert.Equal(t, "aaaacc", res)

	res, err = unpack("4abcd")
	assert.Error(t, err)
	assert.Equal(t, "", res)

	res, err = unpack("a45")
	assert.Error(t, err)
	assert.Equal(t, "", res)

	res, err = unpack("aa2c0")
	assert.Error(t, err)
	assert.Equal(t, "", res)

	res, err = unpack(`abc\3\2`)
	assert.NoError(t, err)
	assert.Equal(t, "abc32", res)

	res, err = unpack(`qwe\45`)
	assert.NoError(t, err)
	assert.Equal(t, "qwe44444", res)

	res, err = unpack(`qwe\\5`)
	assert.NoError(t, err)
	assert.Equal(t, `qwe\\\\\`, res)
}
