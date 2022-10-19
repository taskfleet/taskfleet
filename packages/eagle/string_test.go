package eagle_test

import (
	"encoding/json"
	"errors"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/eagle"
)

type credentials struct {
	Username string
	Password eagle.String
	Test     Foobar
}

type Foobar struct {
	Test string
}

func TestParseString(t *testing.T) {
	data := `{"username": "user", "password": "pass"}`
	var result credentials

	err := json.Unmarshal([]byte(data), &result)
	require.Nil(t, err)

	assert.Equal(t, "user", result.Username)
	assert.Equal(t, "pass", result.Password.Value())
}

func TestParseFile(t *testing.T) {
	data := `{"username": "user", "password": {"file": "testdata/password"}}`
	var result credentials
	err := json.Unmarshal([]byte(data), &result)
	require.Nil(t, err)
	assert.Equal(t, "user", result.Username)
	assert.Equal(t, "securepass", result.Password.Value())
}

func TestParseInvalid(t *testing.T) {
	data := `{"username": "user", "password": {"test": ""}}`
	var result credentials
	err := json.Unmarshal([]byte(data), &result)
	assert.True(t, eagle.IsErrInvalidType(err))
}

func TestParseMissingFile(t *testing.T) {
	data := `{"username": "user", "password": {"file": "testdata/missing"}}`
	var result credentials
	err := json.Unmarshal([]byte(data), &result)
	assert.True(t, errors.Is(err, fs.ErrNotExist))
}
