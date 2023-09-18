package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rapatao/pr-checker-go/domain"
	"github.com/rapatao/pr-checker-go/outgen"
	"github.com/rapatao/pr-checker-go/processor"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

const defaultInstallPath = "/Library/Application Support/xbar/plugins/"

func main() {
	installPath := fmt.Sprintf("%s/%s", os.Getenv("HOME"), defaultInstallPath)

	install := flag.String("install", installPath, "install xbar script")
	flag.Parse()

	runinstall := false
	flag.Visit(func(f *flag.Flag) {
		runinstall = f.Name == "install"
	})

	if runinstall {
		xbarInstall(*install)

		os.Exit(0)
	}

	extract()
}

func extract() {
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

func xbarInstall(path string) {
	log.Printf("installing xbar script at %s", path)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	app := os.Args[0]
	app = strings.TrimLeft(app, "./")

	script := fmt.Sprintf(
		`#!/usr/bin/env sh

if command -v pr-checker-go > /dev/null 2>&1; then
  %s/%s
fi
`, dir, app)

	err = os.WriteFile(path+"/pr-checker.30m.sh", []byte(script), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	println(script)

}
