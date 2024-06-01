package outgen

import (
	"fmt"
	"github.com/rapatao/pr-checker-go/domain"
	"sort"
	"strings"
	"time"
)

type Color string

var (
	NoColor   Color = ""
	OkColor   Color = "| color=#538d22"
	FailColor Color = "| color=#941b0c"
)

type XBarOutGen struct{}

func NewXBarOutGen() domain.OutGen {
	return &XBarOutGen{}
}

func (o *XBarOutGen) Generate(prs []domain.PullRequest) {
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

		fmt.Printf("%s (%d) | href=%s\n", repository, len(prs), prs[0].RepositoryURL)

		for _, pr := range prs {
			prefix := ""
			titleColor := OkColor
			if pr.IsDraft {
				prefix = "(DRAFT) "
				titleColor = NoColor
			}

			if pr.Mergeable != "MERGEABLE" {
				prefix += fmt.Sprintf("(%s) ", pr.Mergeable)
				titleColor = FailColor
			}

			prTitle := strings.ReplaceAll(pr.Title, "|", "ǀ")

			fmt.Printf("-- %s%s | href=%s %s\n", prefix, prTitle, pr.Link, titleColor)
			fmt.Printf("-- issue: #%d by %s\n", pr.Number, pr.Author)
			fmt.Printf("-- created at %v\n", pr.CreatedAt.Format(time.RFC1123))
			fmt.Printf("-- updated at %v\n", pr.UpdatedAt.Format(time.RFC1123))

			if pr.ReviewDecision != "" {
				stateColor := FailColor
				if pr.ReviewDecision == "APPROVED" {
					stateColor = OkColor
				}
				fmt.Printf("-- state: %s %s\n", pr.ReviewDecision, stateColor)
			}

			if pr.CheckStatus != "" {
				checkColor := OkColor
				if pr.CheckStatus == "FAILURE" {
					checkColor = FailColor
				}
				fmt.Printf("-- checks: %s %s\n", pr.CheckStatus, checkColor)
			}

			fmt.Printf("-----\n")
		}
	}

	fmt.Printf("---\n")
	fmt.Printf("Last update: %s\n", time.Now().Format(time.RFC1123))
}

var _ domain.OutGen = (*XBarOutGen)(nil)
