package cli

import "github.com/linuxsuren/http-downloader/pkg/exec"

// Script represents a script object
type Script struct {
	Kind        string
	Title       string
	Content     string
	Dir         string
	KeepScripts bool
	Execer      exec.Execer
}

// ScriptRunner is the interface of a common runner
type ScriptRunner interface {
	Run() error
	GetTitle() string
}

// NewScriptRunners returns the instance of ScriptRunners
func NewScriptRunners() ScriptRunners {
	return []ScriptRunner{&QuitRunner{}}
}

// ScriptRunners is an alias of the ScriptRunner slice
type ScriptRunners []ScriptRunner

// GetTitles returns all the titles
func (s ScriptRunners) GetTitles() (titles []string) {
	titles = make([]string, len(s))
	for i, r := range s {
		titles[i] = r.GetTitle()
	}
	return
}

// GetRunner returns the runner by title
func (s ScriptRunners) GetRunner(title string) ScriptRunner {
	for _, runner := range s {
		if runner.GetTitle() == title {
			return runner
		}
	}
	return nil
}

// Size returns the size of the script runners
func (s ScriptRunners) Size() int {
	return len(s)
}

// QuitRunner represents a runner for quit
type QuitRunner struct{}

// Run does nothing
func (r *QuitRunner) Run() error {
	return nil
}

// GetTitle returns the title
func (r *QuitRunner) GetTitle() string {
	return "Quit"
}
