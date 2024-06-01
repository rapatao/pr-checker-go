package outgen

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/rapatao/pr-checker-go/domain"
	"log"
	"os"
	"text/template"
	"time"
)

//go:embed xbar_stdout.txt.tmpl
var defaultTemplate string

type StdOutGen struct {
	template string
}

func NewStdOutGen() OutGen {
	var tmpl string

	file, err := os.ReadFile(fmt.Sprintf("%s/.pr-checker.tmpl", os.Getenv("HOME")))
	switch {
	case err == nil:
		tmpl = string(file)
	case errors.Is(err, os.ErrNotExist):
		tmpl = defaultTemplate
	default:
		log.Fatal(err)
	}

	return &StdOutGen{
		tmpl,
	}
}

func (o *StdOutGen) Generate(prs []domain.PullRequest) {
	grouped := map[string][]domain.PullRequest{}
	for _, pr := range prs {
		list, ok := grouped[pr.Repository]
		if !ok {
			list = []domain.PullRequest{}
		}
		list = append(list, pr)

		grouped[pr.Repository] = list
	}

	tmpl, err := template.New("xbar").Parse(o.template)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(os.Stdout, map[string]interface{}{
		"Total":       len(prs),
		"Requests":    grouped,
		"GeneratedAt": time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}
}

var _ OutGen = (*StdOutGen)(nil)
