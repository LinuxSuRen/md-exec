package cli

import (
	"context"
	"os"
	"path"
)

// GroovyScript represents the Python script
type GroovyScript struct {
	Script
}

// Run executes the script
func (s *GroovyScript) Run(ctx context.Context) (err error) {
	var shellFile string
	if shellFile, err = writeAsShell(s.Content, s.Dir); err == nil {
		if !s.KeepScripts {
			defer func() {
				_ = os.RemoveAll(shellFile)
			}()
		}

		err = s.Execer.RunCommandInDir("groovy", s.Dir, path.Base(shellFile))
	}
	return
}

// GetTitle returns the title of this script
func (s *GroovyScript) GetTitle() string {
	return s.Title
}

var _ ScriptRunner = &GroovyScript{}
