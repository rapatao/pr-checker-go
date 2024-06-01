package domain

type OutGen interface {
	Generate(prs []PullRequest)
}
