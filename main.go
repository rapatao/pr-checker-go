package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"fyne.io/systray"
	"github.com/rapatao/pr-checker-go/domain"
	"github.com/rapatao/pr-checker-go/processor"
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
	systray.SetTitle("PR")
	systray.SetTooltip("GitHub PR Checker")

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

	systray.ResetMenu()
	mLoading := systray.AddMenuItem("Updating...", "Checking PRs")
	mLoading.Disable()

	home := os.Getenv("HOME")
	file, err := os.ReadFile(fmt.Sprintf("%s/.pr-checker.yml", home))
	if err != nil {
		systray.ResetMenu()
		systray.AddMenuItem("Error reading config", "").Disable()
		addFooter(time.Now())
		return
	}

	var config domain.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		systray.ResetMenu()
		systray.AddMenuItem("Error parsing config", "").Disable()
		addFooter(time.Now())
		return
	}

	prs := processor.Process(context.Background(), &config)
	lastUpdate := time.Now()

	systray.ResetMenu()
	systray.SetTitle(fmt.Sprintf("%d PRs", len(prs)))

	prsByRepo := make(map[string][]domain.PullRequest)
	for _, pr := range prs {
		prsByRepo[pr.Repository] = append(prsByRepo[pr.Repository], pr)
	}

	for repo, repoPrs := range prsByRepo {
		mRepo := systray.AddMenuItem(fmt.Sprintf("%s (%d)", repo, len(repoPrs)), "")
		go func(url string) {
			for range mRepo.ClickedCh {
				openBrowser(url)
			}
		}(repoPrs[0].RepositoryURL)

		for _, pr := range repoPrs {
			m := mRepo.AddSubMenuItem(pr.Title, "")
			go func(url string) {
				for range m.ClickedCh {
					openBrowser(url)
				}
			}(pr.Link)
		}
	}

	addFooter(lastUpdate)
}

func addFooter(lastUpdate time.Time) {
	systray.AddSeparator()
	systray.AddMenuItem(fmt.Sprintf("Last update: %s", lastUpdate.Format("15:04:05")), "").Disable()
	mRefresh := systray.AddMenuItem("Refresh now", "Force update")
	go handleRefreshClick(mRefresh)

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit app")
	go func() {
		for range mQuit.ClickedCh {
			systray.Quit()
		}
	}()
}

func handleRefreshClick(m *systray.MenuItem) {
	for range m.ClickedCh {
		select {
		case refreshChan <- struct{}{}:
		default:
			// Already refreshing, skip
		}
	}
}

func openBrowser(url string) {
	exec.Command("open", url).Run()
}

func onExit() {
	// Cleanup
}
