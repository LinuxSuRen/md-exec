package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/linuxsuren/http-downloader/pkg/exec"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "")
	assert.Nil(t, err)
	defer func() {
		_ = f.Close()
	}()

	tests := []struct {
		name   string
		args   []string
		hasErr bool
	}{{
		name:   "no argument",
		args:   nil,
		hasErr: true,
	}, {
		name:   "a single file which is not Markdown",
		args:   []string{f.Name()},
		hasErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCommand(exec.FakeExecer{}, &bytes.Buffer{})
			assert.True(t, cmd.HasExample())

			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			assert.Equal(t, tt.hasErr, err != nil)
		})
	}
}
