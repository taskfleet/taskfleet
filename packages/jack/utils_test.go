package jack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustNoError(t *testing.T) {
	x := 3
	var err error
	r := Must(x, err)
	assert.Equal(t, x, r)
}

func TestMustError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic but did not record it")
		}
	}()
	x := 3
	err := fmt.Errorf("this is an error")
	Must(x, err)
}
