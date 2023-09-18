package processor

import (
	"context"
	"fmt"
	"github.com/google/go-github/v55/github"
	"github.com/rapatao/pr-checker-go/domain"
	"log"
	"strings"
)

func extractGitHub(ctx context.Context, service *domain.Service) []domain.PullRequest {
	client := github.NewClient(nil).
		WithAuthToken(service.Token)

	var prs []domain.PullRequest

	//filter := "type:pr"
	filter := "type:pr state:open"

	if len(service.Author) > 0 {
		filter = fmt.Sprintf("%s author:%s", filter, service.Author)
	}

	for _, repository := range service.Repositories {
		filter = fmt.Sprintf("%s repo:%s", filter, repository)
	}

	issues, response, err := client.Search.Issues(ctx, filter, &github.SearchOptions{
		Sort: "created",
	})
	if err != nil {
		log.Fatalf("search %s returned %d, %v", filter, response.StatusCode, err)
	}

	for _, issue := range issues.Issues {
		repo := issue.GetHTMLURL()
		if j := strings.LastIndex(repo, "/pull"); j >= 0 {
			repo = repo[:j]
		}

		prs = append(prs, domain.PullRequest{
			Service:    service.Name,
			Repository: repo,
			Title:      issue.GetTitle(),
			Number:     issue.GetNumber(),
			Link:       issue.GetHTMLURL(),
			CreatedAt:  issue.GetCreatedAt().Time,
			UpdatedAt:  issue.GetUpdatedAt().Time,
			Author:     issue.GetUser().GetLogin(),
		})
	}

	return prs
}
