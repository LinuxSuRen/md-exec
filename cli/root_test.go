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

func TestParseMarkdownRunner(t *testing.T) {
	opt := &option{}
	runners, err := opt.parseMarkdownRunner("../README.md")
	if assert.Nil(t, err) {
		assert.True(t, len(runners) > 0)
		assert.NotNil(t, runners.GetRunner("Variable Input Hello World"))
		assert.NotNil(t, runners.GetRunner("Python Hello World"))
		assert.NotNil(t, runners.GetRunner("Run long time"))
		assert.NotNil(t, runners.GetRunner("Golang Hello World"))
	}
}

func TestParseMarkdownRunners(t *testing.T) {
	opt := &option{}
	runners, err := opt.parseMarkdownRunners([]string{"../README.md"})
	if assert.Nil(t, err) {
		assert.True(t, len(runners) > 0)
		assert.NotNil(t, runners.GetRunner("Variable Input Hello World"))
		assert.NotNil(t, runners.GetRunner("Python Hello World"))
		assert.NotNil(t, runners.GetRunner("Run long time"))
		assert.NotNil(t, runners.GetRunner("Golang Hello World"))
	}
}
