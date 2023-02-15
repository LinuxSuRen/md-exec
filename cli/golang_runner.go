package cli

import (
	"fmt"
	"os"
	"os/exec"
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
	if shellFile, err = writeAsShell(fmt.Sprintf(sampleGo, s.Content), s.Dir); err != nil {
		fmt.Println(err)
		return
	}

	goSourceFile := fmt.Sprintf("%s.go", shellFile)
	os.Rename(shellFile, goSourceFile)

	if !s.KeepScripts {
		defer func() {
			_ = os.RemoveAll(goSourceFile)
		}()
	}

	var goExec string
	if goExec, err = exec.LookPath("go"); err != nil {
		return
	}

	cmd := exec.Command(goExec, "run", goSourceFile)
	cmd.Env = os.Environ()

	var output []byte
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(string(output), err)
		return
	}
	fmt.Print(string(output))
	return
}

// GetTitle returns the title of this script
func (s *GolangScript) GetTitle() string {
	return s.Title
}

var _ ScriptRunner = &GolangScript{}
