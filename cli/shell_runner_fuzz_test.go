package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzInputRequest(f *testing.F) {
	f.Add("a", "b")
	f.Fuzz(func(t *testing.T, key, value string) {
		ok, pair, err := isInputRequest(fmt.Sprintf("%s=%s", key, value))
		assert.True(t, ok)
		assert.Equal(t, key, pair[0])
		assert.Equal(t, value, pair[1])
		assert.Nil(t, err)
		if err != nil {
			t.Fail()
		}
	})
}
