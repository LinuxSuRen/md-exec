package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		hasErr bool
	}{{
		name:   "no argument",
		args:   nil,
		hasErr: true,
	}, {
		name:   "more than one argument",
		args:   []string{"a", "b"},
		hasErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCommand()
			assert.True(t, cmd.HasExample())

			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			assert.Equal(t, tt.hasErr, err != nil)
		})
	}
}
