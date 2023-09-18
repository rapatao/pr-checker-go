package executors

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Installer struct {
	plugin string
}

func (i *Installer) Prepare() {
	installPath := fmt.Sprintf("%s/%s",
		os.Getenv("HOME"),
		"Library/Application Support/xbar/plugins/pr-checker.30m.sh")

	flag.StringVar(&i.plugin, "install", installPath, "xbar plugin full path")
}

func (i *Installer) Run() {
	log.Printf("installing xbar script\n")

	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	script := fmt.Sprintf(
		`#!/usr/bin/env sh

if command -v %s > /dev/null 2>&1; then
  %s
fi
`, executable, executable)

	err = os.WriteFile(i.plugin, []byte(script), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Plugin created at: %s\n\n%s", i.plugin, script)
}

func (i *Installer) Is() bool {
	valid := false
	flag.Visit(func(f *flag.Flag) {
		valid = f.Name == "install"
	})

	return valid
}

var _ Executor = (*Installer)(nil)
