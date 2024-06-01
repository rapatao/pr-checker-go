package outgen

import "github.com/rapatao/pr-checker-go/domain"

type OutGen interface {
	Generate(prs []domain.PullRequest)
}
