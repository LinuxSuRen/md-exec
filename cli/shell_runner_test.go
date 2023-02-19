package cli

import (
	"context"
	"testing"

	"github.com/linuxsuren/http-downloader/pkg/exec"

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
			err := shell.Run(context.Background())
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

func TestFindPotentialCommands(t *testing.T) {
	tests := []struct {
		name   string
		cmd    string
		expect []string
	}{{
		name:   "oneline cmd",
		cmd:    "k3d create cluster",
		expect: []string{"k3d"},
	}, {
		name:   "empty",
		cmd:    "",
		expect: nil,
	}, {
		name: "multiple lines",
		cmd: `k3d create cluster
docker ps`,
		expect: []string{"k3d", "docker"},
	}, {
		name:   "with extra whitespace",
		cmd:    " k3d   create    cluster",
		expect: []string{"k3d"},
	}, {
		name: "with EOF",
		cmd: `EOF
k3d create cluster`,
		expect: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findPotentialCommands(tt.cmd)
			assert.Equal(t, tt.expect, result)
		})
	}
}
