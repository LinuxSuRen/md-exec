// Package cli provides all the commands
package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/golang-commonmark/markdown"
	"github.com/linuxsuren/http-downloader/pkg/exec"
	"github.com/spf13/cobra"
)

// should be injected during the build process
var version string

// NewRootCommand returns the instance of cobra.Command
func NewRootCommand(execer exec.Execer, out io.Writer) (cmd *cobra.Command) {
	opt := &option{
		execer: execer,
	}
	cmd = &cobra.Command{
		Use:     "mde",
		Example: "mde README.md",
		Args:    cobra.MinimumNArgs(1),
		RunE:    opt.runE,
	}
	cmd.SetOut(out)
	cmd.Version = version
	flags := cmd.Flags()
	flags.BoolVarP(&opt.loop, "loop", "", true, "Run the Markdown in loop mode.")
	flags.BoolVarP(&opt.keepFilter, "keep-filter", "", true, "Indicate if keep the filter.")
	flags.BoolVarP(&opt.keepScripts, "keep-scripts", "", false, "Indicate if keep the temporary scripts.")
	flags.IntVarP(&opt.pageSize, "page-size", "", 8, "Number of the select items.")
	flags.StringVarP(&opt.filter, "filter", "", "", "Filter of the scripts")
	return
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	var scriptRunners ScriptRunners
	if scriptRunners, err = o.parseMarkdownRunners(args); err == nil && scriptRunners.Size() > 1 {
		for {
			if err = o.executeScripts(scriptRunners); err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
			}

			if !o.loop {
				break
			}
		}
	}
	return
}

func (o *option) parseMarkdownRunners(files []string) (scriptRunners ScriptRunners, err error) {
	scriptRunners = NewScriptRunners()

	for _, mdFilePath := range files {
		var files []string
		if files, err = filepath.Glob(mdFilePath); err == nil {
			for _, file := range files {
				if !strings.HasSuffix(file, ".md") {
					continue
				}
				var runners ScriptRunners
				if runners, err = o.parseMarkdownRunner(file); err != nil {
					break
				}

				scriptRunners = append(scriptRunners, runners...)
			}
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
			Execer:      o.execer,
		}

		switch lang {
		case "shell", "bash", "zsh", "csh", "tcsh", "dash", "fish", "ksh":
			scriptList = append(scriptList, &ShellScript{
				Script:    script,
				ShellType: lang,
			})
		case "python3":
			scriptList = append(scriptList, &PythonScript{
				Script: script,
			})
		case "groovy":
			scriptList = append(scriptList, &GroovyScript{
				Script: script,
			})
		case "golang", "go":
			scriptList = append(scriptList, &GolangScript{
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
	pageSize    int
	filter      string

	execer exec.Execer
}

func (o *option) executeScripts(scriptRunners ScriptRunners) (err error) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	selector := &survey.MultiSelect{
		Message: "Choose the code block to run",
		Options: scriptRunners.GetTitles(),
		VimMode: true,
	}
	var titles []string
	if err = survey.AskOne(selector, &titles,
		survey.WithKeepFilter(o.keepFilter),
		survey.WithFilter(func(filter, value string, index int) (include bool) {
			include = true
			value = strings.ToLower(value)
			if o.filter != "" && value != "quit" {
				include = strings.Contains(value, o.filter)
			}
			return
		}),
		survey.WithPageSize(o.pageSize)); err != nil {
		return
	}

	for _, title := range titles {
		if title == "Quit" {
			o.loop = false
			break
		}

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-c
			cancel()
		}()

		if runner := scriptRunners.GetRunner(title); runner == nil {
			fmt.Println("cannot found runner:", title)
		} else if err = runner.Run(ctx); err != nil {
			break
		}
	}
	return
}
