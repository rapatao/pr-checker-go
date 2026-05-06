package ui

import (
	"fmt"
	"os/exec"
	"sort"
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
	systray.SetTitle("🔄")
	systray.ResetMenu()
	mLoading := systray.AddMenuItem("Updating...", "Checking PRs")
	mLoading.Disable()
	AddFooter(time.Now())
}

func SetError(message string) {
	systray.SetTitle("PR")
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

	repoNames := make([]string, 0, len(prsByRepo))
	for repo := range prsByRepo {
		repoNames = append(repoNames, repo)
	}
	sort.Strings(repoNames)

	for _, repo := range repoNames {
		repoPrs := prsByRepo[repo]
		sort.Slice(repoPrs, func(i, j int) bool {
			return repoPrs[i].Number < repoPrs[j].Number
		})

		mRepo := systray.AddMenuItem(fmt.Sprintf("%s (%d)", repo, len(repoPrs)), "")
		go func(url string) {
			for range mRepo.ClickedCh {
				OpenBrowser(url)
			}
		}(repoPrs[0].RepositoryURL)

		for i, pr := range repoPrs {
			if i > 0 {
				mRepo.AddSeparator()
			}
			m := mRepo.AddSubMenuItem(fmt.Sprintf("#%d: %s", pr.Number, pr.Title), "")
			mRepo.AddSubMenuItem(fmt.Sprintf("   📅 %s", pr.CreatedAt.Format("2006-01-02 15:04")), "").Disable()

			reviewIcon := "👀"
			switch pr.ReviewDecision {
			case "APPROVED":
				reviewIcon = "✅"
			case "CHANGES_REQUESTED":
				reviewIcon = "⚠️"
			}
			mRepo.AddSubMenuItem(fmt.Sprintf("   %s %s", reviewIcon, pr.ReviewDecision), "").Disable()

			statusIcon := "⏳"
			switch pr.CheckStatus {
			case "SUCCESS":
				statusIcon = "✅"
			case "FAILURE", "ERROR":
				statusIcon = "❌"
			}
			mRepo.AddSubMenuItem(fmt.Sprintf("   %s %s", statusIcon, pr.CheckStatus), "").Disable()
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
