package cli

import (
	"github.com/linuxsuren/http-downloader/pkg/exec"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPythonRunner(t *testing.T) {
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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell := &PythonScript{
				Script: Script{
					Title:   tt.title,
					Content: tt.cmd,
					Execer:  exec.FakeExecer{},
				},
			}
			assert.Equal(t, tt.title, shell.GetTitle())
			err := shell.Run()
			assert.Equal(t, tt.hasErr, err != nil)
		})
	}
}
