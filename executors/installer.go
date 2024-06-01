package executors

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

//go:embed install.sh.tmpl
var scriptTemplate string

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
	log.Println("installing xbar script", i.plugin)

	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("install").Parse(scriptTemplate)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(i.plugin, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, map[string]string{
		"executable": executable,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Plugin created at: %s\n", i.plugin)
}

func (i *Installer) Is() bool {
	valid := false
	flag.Visit(func(f *flag.Flag) {
		valid = f.Name == "install"
	})

	return valid
}

var _ Executor = (*Installer)(nil)
