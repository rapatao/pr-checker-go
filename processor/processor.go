package processor

import (
	"context"
	"github.com/rapatao/pr-checker-go/domain"
)

type void struct{}

var nothing void

func Process(ctx context.Context, config *domain.Config) []domain.PullRequest {
	prs := make(map[domain.PullRequest]void)

	for _, service := range config.Services {
		for _, pr := range extractGitHub(ctx, &service) {
			prs[pr] = nothing
		}
	}

	result := make([]domain.PullRequest, 0, len(prs))

	for pr := range prs {
		result = append(result, pr)
	}

	return result
}
