package outgen

import (
	"fmt"
	"github.com/rapatao/pr-checker-go/domain"
	"time"
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

	for repository, prs := range grouped {
		fmt.Printf("%s (%d) | href=%s\n", repository, len(prs), prs[0].RepositoryURL)

		for _, pr := range prs {
			prefix := ""
			titleColor := "#a0db8e"
			if pr.IsDraft {
				prefix = "(DRAFT) "
				titleColor = "#dbdbdb"
			}

			fmt.Printf("-- %s%s | size=14 color=%s href=%s\n", prefix, pr.Title, titleColor, pr.Link)
			fmt.Printf("-- issue: #%d by %s | size=12 color=#aba9bf\n", pr.Number, pr.Author)
			fmt.Printf("-- created at %v | size=12 color=#aba9bf\n", pr.CreatedAt)
			fmt.Printf("-- updated at %v | size=12 color=#aba9bf\n", pr.UpdatedAt)
			fmt.Printf("--\n")
		}
	}

	fmt.Printf("---\n")
	fmt.Printf("Last update: %s\n", time.Now().Format(time.ANSIC))
}
