package processor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rapatao/pr-checker-go/domain"
	"github.com/stretchr/testify/assert"
)

func TestGiteaExtractor_Extract(t *testing.T) {
	routes := map[string]string{
		"/api/v1/repos/issues/search": `[
			{"number": 1, "repository": {"full_name": "org/app"}},
			{"number": 2, "repository": {"full_name": "org/app"}},
			{"number": 3, "repository": {"full_name": "org/other"}}
		]`,
		"/api/v1/repos/org/app/pulls/1": `{
			"number": 1, "title": "Add feature", "html_url": "https://gitea.local/org/app/pulls/1",
			"draft": true, "mergeable": true, "user": {"login": "alice"}, "head": {"sha": "abc"},
			"base": {"repo": {"full_name": "org/app", "html_url": "https://gitea.local/org/app"}}
		}`,
		"/api/v1/repos/org/app/pulls/2":            `{"number": 2, "user": {"login": "bob"}}`,
		"/api/v1/repos/org/app/commits/abc/status": `{"state": "failure"}`,
		"/api/v1/repos/org/app/pulls/1/reviews": `[
			{"state": "REQUEST_CHANGES", "user": {"login": "carol"}},
			{"state": "APPROVED", "user": {"login": "carol"}},
			{"state": "REQUEST_CHANGES", "dismissed": true, "user": {"login": "dave"}}
		]`,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "token secret", r.Header.Get("Authorization"))
		body, ok := routes[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	prs := NewGiteaExtractor().Extract(context.Background(), &domain.Service{
		URL:          server.URL + "/",
		Token:        "secret",
		Author:       "Alice",
		Repositories: []string{"org/app"},
	})

	assert.Equal(t, []domain.PullRequest{{
		Repository:     "org/app",
		RepositoryURL:  "https://gitea.local/org/app",
		Title:          "Add feature",
		Number:         1,
		Link:           "https://gitea.local/org/app/pulls/1",
		Author:         "alice",
		IsDraft:        true,
		ReviewDecision: "APPROVED",
		Mergeable:      "MERGEABLE",
		CheckStatus:    "FAILURE",
	}}, prs)
}
