package processor

import (
	"context"
	"github.com/rapatao/pr-checker-go/domain"
)

type void struct{}

var nothing void

type Extractor interface {
	Extract(ctx context.Context, service *domain.Service) []domain.PullRequest
}
