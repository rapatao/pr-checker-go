package domain

import "time"

type Config struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name         string   `yaml:"name"`
	Provider     string   `yaml:"provider"`
	Token        string   `yaml:"token"`
	Repositories []string `yaml:"repositories"`
}

type PullRequest struct {
	Service    string
	Repository string
	Title      string
	Number     int
	Link       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Author     string
	IsDraft    bool
}
