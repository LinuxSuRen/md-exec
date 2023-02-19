package cli

import (
	"context"
	"os"
	"path"
)

// PythonScript represents the Python script
type PythonScript struct {
	Script
}

// Run executes the script
func (s *PythonScript) Run(ctx context.Context) (err error) {
	var shellFile string
	if shellFile, err = writeAsShell(s.Content, s.Dir); err == nil {
		if !s.KeepScripts {
			defer func() {
				_ = os.RemoveAll(shellFile)
			}()
		}

		err = s.Execer.RunCommandInDir("python3", s.Dir, path.Base(shellFile))
	}
	return
}

// GetTitle returns the title of this script
func (s *PythonScript) GetTitle() string {
	return s.Title
}

var _ ScriptRunner = &PythonScript{}
