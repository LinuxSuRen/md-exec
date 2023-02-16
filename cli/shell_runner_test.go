package cli

import (
	"fmt"
	"github.com/linuxsuren/http-downloader/pkg/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShell(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		shellType string
		cmd       string
		hasErr    bool
	}{{
		name:      "normal bash",
		title:     "title",
		shellType: "bash",
		cmd: `#!title: title
echo 1`,
		hasErr: false,
	}, {
		name:      "shell type is empty",
		title:     "title",
		shellType: "",
		cmd: `#!title: title
echo 1`,
		hasErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell := &ShellScript{
				Script: Script{
					Title:   tt.title,
					Content: tt.cmd,
					Execer:  exec.FakeExecer{},
				},
				ShellType: tt.shellType,
			}
			assert.Equal(t, tt.title, shell.GetTitle())
			err := shell.Run()
			assert.Equal(t, tt.hasErr, err != nil)
		})
	}
}

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
