package executors

import (
	"context"
	"fmt"
	"github.com/rapatao/pr-checker-go/domain"
	"github.com/rapatao/pr-checker-go/outgen"
	"github.com/rapatao/pr-checker-go/processor"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type XBarGen struct{}

func (X *XBarGen) Prepare() {}

func (X *XBarGen) Run() {
	log.SetPrefix("#")

	home := os.Getenv("HOME")
	file, err := os.ReadFile(fmt.Sprintf("%s/.pr-checker.yml", home))
	if err != nil {
		log.Fatal(err)
	}

	var config domain.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pullRequests := processor.Process(ctx, &config)

	outgen.ForXBar(pullRequests)
}

func (X *XBarGen) Is() bool {
	return true
}

var _ Executor = (*XBarGen)(nil)
