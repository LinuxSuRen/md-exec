package cli

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestParseMarkdownRunnerWithCRLF(t *testing.T) {
	// simulate a checkout on Windows where core.autocrlf=true converts
	// every text file to CRLF line endings
	mdContent, err := os.ReadFile("../README.md")
	assert.Nil(t, err)

	lfContent := bytes.ReplaceAll(mdContent, []byte("\r\n"), []byte("\n"))
	crlfContent := bytes.ReplaceAll(lfContent, []byte("\n"), []byte("\r\n"))

	mdFile := filepath.Join(t.TempDir(), "README.md")
	assert.Nil(t, os.WriteFile(mdFile, crlfContent, 0o600))

	opt := &option{}
	runners, err := opt.parseMarkdownRunner(mdFile)
	if assert.Nil(t, err) {
		assert.True(t, len(runners) > 0)
		assert.NotNil(t, runners.GetRunner("Variable Input Hello World"))
		assert.NotNil(t, runners.GetRunner("Python Hello World"))
		assert.NotNil(t, runners.GetRunner("Run long time"))
		assert.NotNil(t, runners.GetRunner("Golang Hello World"))
	}
}

func TestRunWithNoRunnableBlocks(t *testing.T) {
	mdFile := filepath.Join(t.TempDir(), "nothing.md")
	assert.Nil(t, os.WriteFile(mdFile, []byte("# just text\n\nno code blocks here\n"), 0o600))

	buf := new(bytes.Buffer)
	cmd := NewRootCommand(exec.FakeExecer{}, buf)
	cmd.SetArgs([]string{mdFile})
	assert.Nil(t, cmd.Execute())
	assert.Contains(t, buf.String(), "no runnable code blocks found in:")
	assert.Contains(t, buf.String(), "nothing.md")
}
