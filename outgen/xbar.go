package outgen

import (
	"fmt"
	"github.com/rapatao/pr-checker-go/domain"
	"sort"
	"strings"
	"time"
)

var (
	gray   = "#e6e9e6"
	green  = "#538d22"
	red    = "#941b0c"
	yellow = "#ffc300"
)

func ForXBar(prs []domain.PullRequest) {
	fmt.Printf("PR's #%d\n", len(prs))
	fmt.Printf("---\n")

	grouped := map[string][]domain.PullRequest{}
	for _, pr := range prs {
		list, ok := grouped[pr.Repository]
		if !ok {
			list = []domain.PullRequest{}
		}
		list = append(list, pr)

		grouped[pr.Repository] = list
	}

	repos := make([]string, 0, len(grouped))
	for repository, _ := range grouped {
		repos = append(repos, repository)
	}

	sort.Strings(repos)

	for _, repository := range repos {
		prs := grouped[repository]

		fmt.Printf("%s (%d) | color=%s | href=%s\n", repository, len(prs), gray, prs[0].RepositoryURL)

		for _, pr := range prs {
			prefix := ""
			titleColor := green
			if pr.IsDraft {
				prefix = "(DRAFT) "
				titleColor = gray
			}

			if pr.Mergeable != "MERGEABLE" {
				prefix += fmt.Sprintf("(%s) ", pr.Mergeable)
				titleColor = red
			}

			prTitle := strings.ReplaceAll(pr.Title, "|", "ǀ")

			fmt.Printf("-- %s%s | size=14 | color=%s | href=%s | ansi=false \n", prefix, prTitle, titleColor, pr.Link)
			fmt.Printf("-- issue: #%d by %s | size=12 | color=%s\n", pr.Number, pr.Author, gray)
			fmt.Printf("-- created at %v | size=12 | color=%s\n", pr.CreatedAt, gray)
			fmt.Printf("-- updated at %v | size=12 | color=%s\n", pr.UpdatedAt, gray)

			if pr.ReviewDecision != "" {
				stateColor := green
				if pr.ReviewDecision != "APPROVED" {
					stateColor = yellow
				}
				fmt.Printf("-- state: %s | size=12 | color=%s\n", pr.ReviewDecision, stateColor)
			}

			if pr.CheckStatus != "" {
				checkColor := green
				if pr.CheckStatus == "FAILURE" {
					checkColor = red
				}
				fmt.Printf("-- checks: %s | size=12 | color=%s\n", pr.CheckStatus, checkColor)
			}

			fmt.Printf("-----\n")
		}
	}

	fmt.Printf("---\n")
	fmt.Printf("Last update: %s\n", time.Now().Format(time.ANSIC))
}
