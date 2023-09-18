package main

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"pr-checker-go/internal/domain"
	"pr-checker-go/internal/output"
	"pr-checker-go/internal/processor"
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
