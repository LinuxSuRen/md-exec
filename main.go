package main

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
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: mde sample.md")
		return
	}

	mdFilePath := os.Args[1]
	mdFile, err := os.ReadFile(mdFilePath)
	if err != nil {
		panic(err)
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
	execute(cmdMap, contextDir)
}

func execute(cmdMap map[string][]string, contextDir string) (err error) {
	var items []string
	for key := range cmdMap {
		items = append(items, key)
	}

	selector := &survey.MultiSelect{
		Message: "Choose the code block to run",
		Options: items,
	}
	titles := []string{}
	if err = survey.AskOne(selector, &titles); err != nil {
		return
	}

	for _, title := range titles {
		cmds := cmdMap[title]
		for _, cmdLine := range cmds {
			var shellFile string
			if shellFile, err = writeAsShell(cmdLine, contextDir); err != nil {
				fmt.Println(err)
				break
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
				break
			}
			fmt.Print(string(output))
		}
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
