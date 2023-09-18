package executors

import (
	"flag"
	"log"
)

var executors = []Executor{
	&Installer{},
	&XBarGen{},
}

type Executor interface {
	Prepare()
	Run()
	Is() bool
}

func Parse() Executor {
	for _, exec := range executors { // configure all
		exec.Prepare()
	}

	flag.Parse()

	for _, exec := range executors { // return first match
		if exec.Is() {
			return exec
		}
	}

	log.Fatal("none executor matches")

	return nil
}
