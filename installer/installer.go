package installer

import (
	_ "embed"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed launchagent.plist.tmpl
var plistTemplate string

func Install() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		return err
	}

	home := os.Getenv("HOME")
	plistPath := filepath.Join(home, "Library", "LaunchAgents", "com.rapatao.pr-checker.plist")

	f, err := os.Create(plistPath)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, map[string]string{"Path": absPath})
}
