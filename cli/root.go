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
	"github.com/golang-commonmark/markdown"
	"github.com/spf13/cobra"
)

// should be inject during the build process
var version string

func NewRootCommand() (cmd *cobra.Command) {
	opt := &option{}
	cmd = &cobra.Command{
		Use:     "mde",
		Example: "mde README.md",
		Args:    cobra.ExactArgs(1),
		RunE:    opt.runE,
	}
	cmd.Version = version
	flags := cmd.Flags()
	flags.BoolVarP(&opt.loop, "loop", "", true, "Run the Markdown in loop mode.")
	flags.BoolVarP(&opt.keepFilter, "keep-filter", "", true, "Indicate if keep the filter.")
	return
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	mdFilePath := args[0]

	for {
		err = o.runMarkdown(mdFilePath)

		if !o.loop {
			break
		}
	}
	return
}

func (o *option) runMarkdown(mdFilePath string) (err error) {
	var mdFile []byte
	mdFile, err = os.ReadFile(mdFilePath)
	if err != nil {
		return
	}

	//Parse the markdown
	md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
	tokens := md.Parse(mdFile)

	cmdMap := map[string][]string{}

	// Print the result
	var title string
	for _, t := range tokens {
		var (
			content string
			lang    string
		)

		switch tok := t.(type) {
		case *markdown.Fence:
			content = strings.TrimSpace(tok.Content)
			lang = tok.Params
		}

		if content != "" && lang == "shell" {
			// handle the break line
			breakline := regexp.MustCompile(`\\\n`)
			content = breakline.ReplaceAllString(content, "")

			whitespaces := regexp.MustCompile(` +`)
			content = whitespaces.ReplaceAllString(content, " ")

			lines := strings.Split(content, "\n")
			if len(lines) < 2 {
				continue
			}
			title = lines[0]
			if !strings.HasPrefix(title, "#!title: ") {
				continue
			}
			title = strings.TrimPrefix(title, "#!title: ")
			cmdMap[title] = append(cmdMap[title], lines[1:]...)
		}
	}

	contextDir := path.Dir(mdFilePath)
	// TODO this should be a treemap instead of hashmap
	err = o.execute(cmdMap, contextDir)
	return
}

type option struct {
	loop       bool
	keepFilter bool
}

func (o *option) execute(cmdMap map[string][]string, contextDir string) (err error) {
	var items []string
	for key := range cmdMap {
		items = append(items, key)
	}

	items = append(items, "Quit")
	selector := &survey.MultiSelect{
		Message: "Choose the code block to run",
		Options: items,
	}
	titles := []string{}
	if err = survey.AskOne(selector, &titles, survey.WithKeepFilter(o.keepFilter)); err != nil {
		return
	}

	for _, title := range titles {
		if title == "Quit" {
			o.loop = false
			break
		}
		preDefinedEnv := os.Environ()
		cmds := cmdMap[title]
		for _, cmdLine := range cmds {
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

			err = runCmdLine(cmdLine, contextDir)
			if err != nil {
				break
			}
		}

		// reset the env
		os.Clearenv()
		for _, pair := range preDefinedEnv {
			os.Setenv(strings.Split(pair, "=")[0], strings.Split(pair, "=")[1])
		}
	}
	return
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

func runCmdLine(cmdLine, contextDir string) (err error) {
	var shellFile string
	if shellFile, err = writeAsShell(cmdLine, contextDir); err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = os.RemoveAll(shellFile)
	}()

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

func runAsInlineCommand(cmdLine, contextDir string) (err error) {
	args := strings.Split(cmdLine, " ")
	cmd := strings.TrimSpace(args[0])
	if cmd, err = exec.LookPath(cmd); err != nil {
		err = fmt.Errorf("failed to find '%s'", cmd)
		return
	}

	fmt.Printf("start to run: %s %v\n", cmd, args[1:])
	var output []byte
	cmdRun := exec.Command(cmd, args[1:]...)
	cmdRun.Dir = contextDir
	cmdRun.Env = os.Environ()
	if output, err = cmdRun.CombinedOutput(); err == nil {
		fmt.Print(string(output))
	}
	return
}
