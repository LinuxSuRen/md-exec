// Package main is the entrypoint of this CLI project
package main

import (
	"context"
	"os"

	"github.com/linuxsuren/http-downloader/pkg/exec"
	"github.com/linuxsuren/md-exec/cli"
)

func main() {
	cmd := cli.NewRootCommand(exec.DefaultExecer{}, os.Stdout)
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
