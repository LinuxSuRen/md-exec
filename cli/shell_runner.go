package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// ShellScript represents the shell script
type ShellScript struct {
	Script
}

// Run executes the script
func (s *ShellScript) Run() (err error) {
	// handle the break line
	breakline := regexp.MustCompile(`\\\n`)
	s.Content = breakline.ReplaceAllString(s.Content, "")

	whitespaces := regexp.MustCompile(` +`)
	s.Content = whitespaces.ReplaceAllString(s.Content, " ")

	lines := strings.Split(s.Content, "\n")[1:]

	preDefinedEnv := os.Environ()
	for _, cmdLine := range lines {
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
			os.Setenv(pair[0], pair[1])
			continue
		}

		err = runCmdLine(cmdLine, s.Dir, s.KeepScripts)
		if err != nil {
			break
		}
	}

	// reset the env
	os.Clearenv()
	for _, pair := range preDefinedEnv {
		os.Setenv(strings.Split(pair, "=")[0], strings.Split(pair, "=")[1])
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

func runCmdLine(cmdLine, contextDir string, keepScripts bool) (err error) {
	var shellFile string
	if shellFile, err = writeAsShell(cmdLine, contextDir); err != nil {
		fmt.Println(err)
		return
	}
	if !keepScripts {
		defer func() {
			_ = os.RemoveAll(shellFile)
		}()
	}

	cmd := exec.Command("bash", path.Base(shellFile))
	cmd.Dir = contextDir
	cmd.Env = os.Environ()

	var output []byte
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(string(output), err)
		return
	}
	fmt.Print(string(output))
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