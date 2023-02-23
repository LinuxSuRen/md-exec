package cli

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
)

// GolangScript represents the Golang script
type GolangScript struct {
	Script
}

var sampleGo = `
package main
%s
func main(){
	%s
}
`

// Run executes the script
func (s *GolangScript) Run(ctx context.Context) (err error) {
	s.Content = strings.ReplaceAll(s.Content, "#!title: "+s.Title, "")

	imports := findImports(s.Content)
	body := strings.ReplaceAll(s.Content, imports, "")

	var shellFile string
	if shellFile, err = writeAsShell(fmt.Sprintf(sampleGo, imports, body), s.Dir); err == nil {
		goSourceFile := fmt.Sprintf("%s.go", shellFile)
		os.Rename(shellFile, goSourceFile)

		if !s.KeepScripts {
			defer func() {
				_ = os.RemoveAll(goSourceFile)
			}()
		}

		if err = s.Execer.RunCommandInDir("go", s.Dir, "mod", "init", "github.com/linuxsuren/test"); err != nil {
			return
		}

		if err = s.Execer.RunCommandInDir("go", s.Dir, "mod", "tidy"); err != nil {
			return
		}

		err = s.Execer.RunCommandInDir("go", s.Dir, "run", path.Base(goSourceFile))
	}
	return
}

func findImports(content string) (imports string) {
	for _, line := range strings.Split(content, "\n") {
		if !strings.HasPrefix(line, "import ") {
			continue
		}
		imports += line + "\n"
	}
	imports = strings.TrimSpace(imports)
	return
}

// GetTitle returns the title of this script
func (s *GolangScript) GetTitle() string {
	return s.Title
}

var _ ScriptRunner = &GolangScript{}
