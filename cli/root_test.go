package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInputRequest(t *testing.T) {
	tests := []struct {
		name       string
		cmdLine    string
		expectOK   bool
		expectPair []string
		expectErr  bool
	}{{
		name:       "normal",
		cmdLine:    "name=linuxsuren",
		expectOK:   true,
		expectPair: []string{"name", "linuxsuren"},
		expectErr:  false,
	}, {
		name:       "abnormal variable expression - with extra whitespace",
		cmdLine:    "name = linuxsuren",
		expectOK:   false,
		expectPair: nil,
		expectErr:  false,
	}, {
		name:       "complex characters in pair",
		cmdLine:    "vm=i-dy87owjl",
		expectOK:   true,
		expectPair: []string{"vm", "i-dy87owjl"},
		expectErr:  false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, pair, err := isInputRequest(tt.cmdLine)
			assert.Equal(t, tt.expectOK, ok)
			assert.Equal(t, tt.expectPair, pair)
			assert.Equal(t, tt.expectErr, err != nil)
		})
	}
}

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
