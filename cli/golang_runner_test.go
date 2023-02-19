package cli

import (
	"context"
	"testing"

	"github.com/linuxsuren/http-downloader/pkg/exec"
	"github.com/stretchr/testify/assert"
)

func TestGolangRunner(t *testing.T) {
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
			shell := &GolangScript{
				Script: Script{
					Title:   tt.title,
					Content: tt.cmd,
					Execer:  exec.FakeExecer{},
				},
			}
			assert.Equal(t, tt.title, shell.GetTitle())
			err := shell.Run(context.Background())
			assert.Equal(t, tt.hasErr, err != nil)
		})
	}
}
