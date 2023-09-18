package processor

import (
	"context"
	"github.com/google/go-github/v55/github"
	"github.com/rapatao/pr-checker-go/internal/domain"
	"log"
	"strings"
)

func extractGitHub(ctx context.Context, service *domain.Service) []domain.PullRequest {
	client := github.NewClient(nil).
		WithAuthToken(service.Token)

	var prs []domain.PullRequest

	for _, repository := range service.Repositories {
		split := strings.Split(repository, "/")
		if len(split) != 2 {
			log.Printf("ignoring %s due to malformatting entry \n", repository)

			continue
		}

		opts := &github.PullRequestListOptions{
			Sort: "created",
		}

		list, response, err := client.PullRequests.List(ctx, split[0], split[1], opts)
		if err != nil {
			log.Printf("%s returned %d. %v \n", repository, response.StatusCode, err)

			continue
		}

		for _, pullRequest := range list {

			prs = append(prs, domain.PullRequest{
				Service:    service.Name,
				Repository: repository,
				Title:      pullRequest.GetTitle(),
				Number:     pullRequest.GetNumber(),
				Link:       pullRequest.GetLinks().GetHTML().GetHRef(),
				CreatedAt:  pullRequest.GetCreatedAt().Time,
				UpdatedAt:  pullRequest.GetUpdatedAt().Time,
				Author:     pullRequest.GetUser().GetLogin(),
				IsDraft:    pullRequest.GetDraft(),
			})
		}
	}

	return prs
}
