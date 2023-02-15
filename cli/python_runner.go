package cli

import (
	"fmt"
	"os"
	"os/exec"
)

// PythonScript represents the Python script
type PythonScript struct {
	Script
}

// Run executes the script
func (s *PythonScript) Run() (err error) {
	var shellFile string
	if shellFile, err = writeAsShell(s.Content, s.Dir); err != nil {
		fmt.Println(err)
		return
	}
	if !s.KeepScripts {
		defer func() {
			_ = os.RemoveAll(shellFile)
		}()
	}

	var pyExec string
	if pyExec, err = exec.LookPath("python3"); err != nil {
		return
	}

	cmd := exec.Command(pyExec, shellFile)
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
func (s *PythonScript) GetTitle() string {
	return s.Title
}

var _ ScriptRunner = &PythonScript{}
