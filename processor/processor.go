package processor

import (
	"context"
	"github.com/rapatao/pr-checker-go/domain"
	"log"
	"strings"
)

var extractors = map[string]Extractor{
	"github": NewGitHubExtractor(),
}

func Process(ctx context.Context, config *domain.Config) []domain.PullRequest {
	prs := make(map[domain.PullRequest]void)

	for _, service := range config.Services {
		extractor, ok := extractors[strings.ToLower(service.Provider)]
		if !ok {
			log.Fatalf("service %s is not supported", service.Provider)
		}

		for _, pr := range extractor.Extract(ctx, &service) {
			prs[pr] = nothing
		}
	}

	result := make([]domain.PullRequest, 0, len(prs))

	for pr := range prs {
		result = append(result, pr)
	}

	return result
}
