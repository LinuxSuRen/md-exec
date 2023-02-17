package cli

import (
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
import "fmt"
func main(){
	%s
}
`

// Run executes the script
func (s *GolangScript) Run() (err error) {
	s.Content = strings.ReplaceAll(s.Content, "#!title: "+s.Title, "")

	var shellFile string
	if shellFile, err = writeAsShell(fmt.Sprintf(sampleGo, s.Content), s.Dir); err == nil {
		goSourceFile := fmt.Sprintf("%s.go", shellFile)
		os.Rename(shellFile, goSourceFile)

		if !s.KeepScripts {
			defer func() {
				_ = os.RemoveAll(goSourceFile)
			}()
		}

		err = s.Execer.RunCommandInDir("go", s.Dir, "run", path.Base(goSourceFile))
	}
	return
}

// GetTitle returns the title of this script
func (s *GolangScript) GetTitle() string {
	return s.Title
}

var _ ScriptRunner = &GolangScript{}
