package main

import (
	"github.com/rapatao/pr-checker-go/executors"
)

func main() {
	executor := executors.Parse()
	executor.Run()
}
