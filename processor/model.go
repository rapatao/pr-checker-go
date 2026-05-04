package processor

import (
	"context"

	"github.com/rapatao/pr-checker-go/domain"
)

type Extractor interface {
	Extract(ctx context.Context, service *domain.Service) []domain.PullRequest
}
