package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/golang-commonmark/markdown"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: md-exec sample.md")
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

	//Print the result
	for _, t := range tokens {
		snippet := getSnippet(t)
		snippet.content = strings.TrimSpace(snippet.content)
		if snippet.content != "" {
			if snippet.lang == "shell" {
				lines := strings.Split(snippet.content, "\n")
				for _, line := range lines {
					items := strings.Split(line, " ")
					data, _ := exec.Command(items[0], items[1:]...).CombinedOutput()
					fmt.Print(string(data))
				}
			}
		}
	}
}

// snippet represents the snippet we will output.
type snippet struct {
	content string
	lang    string
}

// getSnippet extract only code snippet from markdown object.
func getSnippet(tok markdown.Token) snippet {
	switch tok := tok.(type) {
	case *markdown.CodeBlock:
		return snippet{
			tok.Content,
			"code",
		}
	case *markdown.CodeInline:
		return snippet{
			tok.Content,
			"code inline",
		}
	case *markdown.Fence:
		return snippet{
			tok.Content,
			tok.Params,
		}
	}
	return snippet{}
}
