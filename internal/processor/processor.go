package processor

import (
	"context"
	"github.com/rapatao/pr-checker-go/internal/domain"
)

func Process(ctx context.Context, config *domain.Config) []domain.PullRequest {
	var prs []domain.PullRequest

	for _, service := range config.Services {
		for _, pr := range extractGitHub(ctx, &service) {
			prs = append(prs, pr)
		}
	}

	return prs
}
