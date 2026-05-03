package processor

import (
	"context"
	"log"
	"strings"

	"github.com/rapatao/pr-checker-go/domain"
)

var extractors = map[string]Extractor{
	"github": NewGitHubExtractor(),
}

func Process(ctx context.Context, config *domain.Config) []domain.PullRequest {
	return process(ctx, config, extractors)
}

func process(ctx context.Context, config *domain.Config, extractors map[string]Extractor) []domain.PullRequest {
	prs := make(map[string]domain.PullRequest)

	for _, service := range config.Services {
		extractor, ok := extractors[strings.ToLower(service.Provider)]
		if !ok {
			log.Fatalf("service %s is not supported", service.Provider)
		}

		for _, pr := range extractor.Extract(ctx, &service) {
			prs[pr.Link] = pr
		}
	}

	result := make([]domain.PullRequest, 0, len(prs))

	for _, pr := range prs {
		result = append(result, pr)
	}

	return result
}
