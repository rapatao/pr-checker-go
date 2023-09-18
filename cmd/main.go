package main

import (
	"context"
	"fmt"
	"github.com/rapatao/pr-checker-go/internal/domain"
	"github.com/rapatao/pr-checker-go/internal/output"
	"github.com/rapatao/pr-checker-go/internal/processor"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func main() {
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

	output.ForXBar(pullRequests)
}
