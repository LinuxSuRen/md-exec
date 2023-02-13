package main

import (
	"fmt"
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
			fmt.Println("start to run:", cmdLine)

			args := strings.Split(cmdLine, " ")
			cmd := strings.TrimSpace(args[0])
			if cmd, err = exec.LookPath(cmd); err != nil {
				fmt.Println("failed to find", cmd)
				continue
			}

			var output []byte
			cmdRun := exec.Command(cmd, args[1:]...)
			cmdRun.Dir = contextDir
			cmdRun.Env = os.Environ()
			if output, err = cmdRun.CombinedOutput(); err == nil {
				fmt.Print(string(output))
			} else {
				fmt.Println(string(output), err)
			}
		}
	}
	return
}
