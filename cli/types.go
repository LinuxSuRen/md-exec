package cli

// Script represents a script object
type Script struct {
	Kind        string
	Title       string
	Content     string
	Dir         string
	KeepScripts bool
}

type ScriptRunner interface {
	Run() error
	GetTitle() string
}

func NewScriptRunners() ScriptRunners {
	return []ScriptRunner{&QuitRunner{}}
}

type ScriptRunners []ScriptRunner

func (s ScriptRunners) GetTitles() (titles []string) {
	titles = make([]string, len(s))
	for i, r := range s {
		titles[i] = r.GetTitle()
	}
	return
}
func (s ScriptRunners) GetRunner(title string) ScriptRunner {
	for _, runner := range s {
		if runner.GetTitle() == title {
			return runner
		}
	}
	return nil
}

type QuitRunner struct{}

func (r *QuitRunner) Run() error {
	return nil
}
func (r *QuitRunner) GetTitle() string {
	return "Quit"
}
