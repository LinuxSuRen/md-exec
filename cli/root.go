// Package cli provides all the commands
package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/golang-commonmark/markdown"
	"github.com/spf13/cobra"
)

// should be inject during the build process
var version string

// NewRootCommand returns the instance of cobra.Command
func NewRootCommand() (cmd *cobra.Command) {
	opt := &option{}
	cmd = &cobra.Command{
		Use:     "mde",
		Example: "mde README.md",
		Args:    cobra.MinimumNArgs(1),
		RunE:    opt.runE,
	}
	cmd.Version = version
	flags := cmd.Flags()
	flags.BoolVarP(&opt.loop, "loop", "", true, "Run the Markdown in loop mode.")
	flags.BoolVarP(&opt.keepFilter, "keep-filter", "", true, "Indicate if keep the filter.")
	flags.BoolVarP(&opt.keepScripts, "keep-scripts", "", false, "Indicate if keep the temporary scripts.")
	return
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	scriptRunners := NewScriptRunners()

	for _, mdFilePath := range args {
		var files []string
		if files, err = filepath.Glob(mdFilePath); err != nil {
			return
		}

		for _, file := range files {
			var runners ScriptRunners
			if runners, err = o.parseMarkdownRunner(file); err != nil {
				break
			}

			scriptRunners = append(scriptRunners, runners...)
		}
	}

	for {
		err = o.executeScripts(scriptRunners)

		if !o.loop {
			break
		}
	}
	return
}

func (o *option) parseMarkdownRunner(mdFilePath string) (scriptList ScriptRunners, err error) {
	var mdFile []byte
	mdFile, err = os.ReadFile(mdFilePath)
	if err != nil {
		return
	}

	// Parse the markdown
	scriptList = ScriptRunners{}
	md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
	tokens := md.Parse(mdFile)

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

		if content == "" {
			continue
		}

		originalContent := content
		lines := strings.Split(content, "\n")
		if len(lines) < 2 {
			continue
		}
		title = lines[0]
		if !strings.HasPrefix(title, "#!title: ") {
			continue
		}
		title = strings.TrimPrefix(title, "#!title: ")

		script := Script{
			Kind:        lang,
			Title:       title,
			Content:     originalContent,
			Dir:         path.Dir(mdFilePath),
			KeepScripts: o.keepScripts,
		}

		switch lang {
		case "shell", "bash":
			scriptList = append(scriptList, &ShellScript{
				Script: script,
			})
		case "python3":
			scriptList = append(scriptList, &PythonScript{
				Script: script,
			})
		}
	}
	return
}

type option struct {
	loop        bool
	keepFilter  bool
	keepScripts bool
}

func (o *option) executeScripts(scriptRunners ScriptRunners) (err error) {
	selector := &survey.MultiSelect{
		Message: "Choose the code block to run",
		Options: scriptRunners.GetTitles(),
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

		if runner := scriptRunners.GetRunner(title); runner == nil {
			fmt.Println("cannot found runner:", title)
		} else if err = runner.Run(); err != nil {
			break
		}
	}
	return
}
