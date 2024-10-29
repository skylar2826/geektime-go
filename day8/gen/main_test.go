package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_Gen(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := gen(buffer, "testdata/user.go")
	require.NoError(t, err)
	assert.Equal(t, ``, buffer.String())
}

func Test_Gen_File(t *testing.T) {
	f, err := os.Create("testdata/user_gen.go")
	err = gen(f, "testdata/user.go")
	require.NoError(t, err)
}
