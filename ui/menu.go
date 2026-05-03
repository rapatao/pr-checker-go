package ui

import (
	"fmt"
	"os/exec"
	"time"

	"fyne.io/systray"
	"github.com/rapatao/pr-checker-go/domain"
)

var refreshCallback func()

func Init(onRefresh func()) {
	refreshCallback = onRefresh
	systray.SetTitle("PR")
	systray.SetTooltip("GitHub PR Checker")
}

func OpenBrowser(url string) {
	if err := exec.Command("open", url).Run(); err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
	}
}

func SetLoading() {
	systray.ResetMenu()
	mLoading := systray.AddMenuItem("Updating...", "Checking PRs")
	mLoading.Disable()
}

func SetError(message string) {
	systray.ResetMenu()
	systray.AddMenuItem(message, "").Disable()
	AddFooter(time.Now())
}

func RenderPRs(prs []domain.PullRequest) {
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
				OpenBrowser(url)
			}
		}(repoPrs[0].RepositoryURL)

		for _, pr := range repoPrs {
			m := mRepo.AddSubMenuItem(pr.Title, "")
			go func(url string) {
				for range m.ClickedCh {
					OpenBrowser(url)
				}
			}(pr.Link)
		}
	}
	AddFooter(time.Now())
}

func AddFooter(lastUpdate time.Time) {
	systray.AddSeparator()
	systray.AddMenuItem(fmt.Sprintf("Last update: %s", lastUpdate.Format("15:04:05")), "").Disable()
	mRefresh := systray.AddMenuItem("Refresh now", "Force update")
	go func() {
		for range mRefresh.ClickedCh {
			refreshCallback()
		}
	}()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit app")
	go func() {
		for range mQuit.ClickedCh {
			systray.Quit()
		}
	}()
}
