package cli

import (
	"context"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/linuxsuren/http-downloader/pkg/installer"

	"github.com/AlecAivazis/survey/v2"
)

// ShellScript represents the shell script
type ShellScript struct {
	Script
	ShellType string
}

// Run executes the script
func (s *ShellScript) Run(ctx context.Context) (err error) {
	s.Content = strings.ReplaceAll(s.Content, "\r\n", "\n")

	lines := strings.Split(s.Content, "\n")[1:]

	preDefinedEnv := os.Environ()
	for i, cmdLine := range lines {
		var pair []string
		var ok bool
		ok, pair, err = isInputRequest(cmdLine)
		if err != nil {
			break
		}

		if ok {
			if pair, err = inputRequest(pair); err != nil {
				break
			}
			_ = os.Setenv(pair[0], pair[1])
			continue
		} else {
			err = s.runCmdLine(strings.Join(lines[i:], "\n"), s.Dir, s.KeepScripts)
			break
		}
	}

	// reset the env
	os.Clearenv()
	for _, pair := range preDefinedEnv {
		_ = os.Setenv(strings.Split(pair, "=")[0], strings.Split(pair, "=")[1])
	}
	return
}

func (s *ShellScript) runCmdLine(cmdLine, contextDir string, keepScripts bool) (err error) {
	var shellFile string
	if shellFile, err = writeAsShell(cmdLine, contextDir); err == nil {
		if !keepScripts {
			defer func() {
				_ = os.RemoveAll(shellFile)
			}()
		}

		if s.ShellType == "shell" || s.ShellType == "" {
			s.ShellType = "bash"
		}

		deps := map[string]string{
			s.ShellType: s.ShellType,
		}
		for _, cmd := range findPotentialCommands(cmdLine) {
			deps[cmd] = cmd
		}

		is := installer.Installer{
			Provider: "github",
		}
		if err = is.CheckDepAndInstall(deps); err == nil {
			err = s.Execer.RunCommandInDir(s.ShellType, contextDir, path.Base(shellFile))
		}
	}
	return
}

func findPotentialCommands(cmdLine string) (cmds []string) {
	// TODO should find a better way to skip EOF part
	if strings.Contains(cmdLine, "EOF") {
		return
	}

	lines := strings.Split(cmdLine, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		reg, err := regexp.Compile(`^(hd|curl|wget|k3d|docker)`)
		if err == nil && reg.Match([]byte(line)) {
			if cmd := reg.Find([]byte(line)); cmd != nil {
				cmds = append(cmds, string(cmd))
			}
		}
	}
	return
}

// GetTitle returns the title of this script
func (s *ShellScript) GetTitle() string {
	return s.Title
}

func isInputRequest(cmdLine string) (ok bool, pair []string, err error) {
	var reg *regexp.Regexp
	if reg, err = regexp.Compile(`^\w+=.+$`); err == nil {
		items := strings.Split(cmdLine, "=")
		if reg.MatchString(cmdLine) && len(items) == 2 {
			pair = []string{strings.TrimSpace(items[0]), strings.TrimSpace(items[1])}
			ok = true
		}
	}
	return
}

func inputRequest(pair []string) (result []string, err error) {
	input := survey.Input{
		Message: pair[0],
		Default: pair[1],
	}
	result = pair

	var value string
	if err = survey.AskOne(&input, &value); err == nil {
		result[1] = value
	}

	return
}

func writeAsShell(content, dir string) (targetPath string, err error) {
	var f *os.File
	if f, err = os.CreateTemp(dir, "sh"); err == nil {
		defer func() {
			_ = f.Close()
		}()

		targetPath = f.Name()
		_, err = io.WriteString(f, content)
	}
	return
}

var _ ScriptRunner = &ShellScript{}
