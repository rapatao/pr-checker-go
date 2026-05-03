package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/rapatao/pr-checker-go/domain"
	"github.com/rapatao/pr-checker-go/processor"
	"github.com/rapatao/pr-checker-go/ui"
	"gopkg.in/yaml.v3"
)

var (
	refreshChan  = make(chan struct{}, 1)
	refreshMutex sync.Mutex
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	ui.Init(triggerRefresh)

	// Trigger initial refresh
	go func() {
		refreshChan <- struct{}{}
	}()

	go updateMenu()
}

func updateMenu() {
	ticker := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-ticker.C:
			refresh()
		case <-refreshChan:
			refresh()
		}
	}
}

func refresh() {
	refreshMutex.Lock()
	defer refreshMutex.Unlock()

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

	prs := processor.Process(context.Background(), &config)
	ui.RenderPRs(prs)
}

func triggerRefresh() {
	select {
	case refreshChan <- struct{}{}:
	default:
		// Already refreshing, skip
	}
}

func onExit() {
	// Cleanup
}
