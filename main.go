package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/rapatao/pr-checker-go/domain"
	"github.com/rapatao/pr-checker-go/installer"
	"github.com/rapatao/pr-checker-go/processor"
	"github.com/rapatao/pr-checker-go/ui"
	"gopkg.in/yaml.v3"
)

var (
	cancelFunc context.CancelFunc
	mu         sync.Mutex
)

func main() {
	installFlag := flag.Bool("install", false, "Install as a launch agent")
	uninstallFlag := flag.Bool("uninstall", false, "Uninstall the launch agent")
	flag.Parse()

	if *installFlag {
		if err := installer.Install(); err != nil {
			fmt.Printf("Failed to install: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Installed successfully!")
		os.Exit(0)
	}

	if *uninstallFlag {
		if err := installer.Uninstall(); err != nil {
			fmt.Printf("Failed to uninstall: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Uninstalled successfully!")
		os.Exit(0)
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	ui.Init(triggerRefresh)

	// Trigger initial refresh
	triggerRefresh()

	go updateMenu()
}

func updateMenu() {
	ticker := time.NewTicker(30 * time.Minute)
	for range ticker.C {
		triggerRefresh()
	}
}

func refresh(ctx context.Context) {
	ui.SetLoading()

	home := os.Getenv("HOME")
	file, err := os.ReadFile(fmt.Sprintf("%s/.pr-checker.yml", home))
	if err != nil {
		ui.SetError("Error reading config")
		return
	}

	var config domain.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		ui.SetError("Error parsing config")
		return
	}

	prs := processor.Process(ctx, &config)

	select {
	case <-ctx.Done():
		return
	default:
		ui.RenderPRs(prs)
	}
}

func triggerRefresh() {
	mu.Lock()
	defer mu.Unlock()

	if cancelFunc != nil {
		cancelFunc()
	}

	var ctx context.Context
	ctx, cancelFunc = context.WithCancel(context.Background())

	go refresh(ctx)
}

func onExit() {
	// Cleanup
}
