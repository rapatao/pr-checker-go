package domain

import "time"

type PullRequest struct {
	Repository     string
	RepositoryURL  string
	Title          string
	Number         int
	Link           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Author         string
	IsDraft        bool
	ReviewDecision string
	Mergeable      string
	CheckStatus    string
}
