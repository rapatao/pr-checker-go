package processor

import (
	"context"
	"testing"
	"time"

	"github.com/rapatao/pr-checker-go/domain"
	"github.com/stretchr/testify/assert"
)

type mockExtractor struct {
	prs []domain.PullRequest
}

func (m *mockExtractor) Extract(ctx context.Context, service *domain.Service) []domain.PullRequest {
	return m.prs
}

func TestProcess_Deduplication(t *testing.T) {
	mock := &mockExtractor{
		prs: []domain.PullRequest{
			{Link: "http://example.com/1", Title: "PR 1"},
			{Link: "http://example.com/1", Title: "PR 1 (Duplicate)"},
		},
	}
	extractors := map[string]Extractor{
		"mock": mock,
	}
	config := &domain.Config{
		Services: []domain.Service{{Provider: "mock"}},
	}

	prs := process(context.Background(), config, extractors)

	assert.Len(t, prs, 1)
	assert.Equal(t, "PR 1 (Duplicate)", prs[0].Title)
}

func TestProcess_EmptyConfig(t *testing.T) {
	config := &domain.Config{Services: []domain.Service{}}
	prs := process(context.Background(), config, make(map[string]Extractor))
	assert.Empty(t, prs)
}

func TestPullRequest_StructFields(t *testing.T) {
	pr := domain.PullRequest{
		Repository: "user/repo",
		Title:      "Test PR",
		Link:       "http://example.com/1",
		CreatedAt:  time.Now(),
		Author:     "user",
	}

	assert.Equal(t, "user/repo", pr.Repository)
	assert.Equal(t, "Test PR", pr.Title)
	assert.Equal(t, "http://example.com/1", pr.Link)
	assert.Equal(t, "user", pr.Author)
}
