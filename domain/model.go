package domain

import "time"

type Config struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name         string   `yaml:"name"`
	Provider     string   `yaml:"provider"`
	Token        string   `yaml:"token"`
	Author       string   `yaml:"author"`
	Repositories []string `yaml:"repositories"`
}

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
}
