package installer

import (
	_ "embed"
	"os"
	"os/exec"
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

	if err := tmpl.Execute(f, map[string]string{"Path": absPath}); err != nil {
		return err
	}

	// Load the agent
	return exec.Command("launchctl", "load", plistPath).Run()
}
