package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/rapatao/pr-checker-go/domain"
)

const giteaPageSize = 50

type GiteaExtractor struct {
	client *http.Client
}

func NewGiteaExtractor() Extractor {
	return &GiteaExtractor{client: http.DefaultClient}
}

type giteaIssue struct {
	Number     int `json:"number"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type giteaPull struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	HTMLURL   string    `json:"html_url"`
	Draft     bool      `json:"draft"`
	Mergeable bool      `json:"mergeable"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
	Head struct {
		Sha string `json:"sha"`
	} `json:"head"`
	Base struct {
		Repo struct {
			FullName string `json:"full_name"`
			HTMLURL  string `json:"html_url"`
		} `json:"repo"`
	} `json:"base"`
}

type giteaReview struct {
	State     string `json:"state"`
	Dismissed bool   `json:"dismissed"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
}

func (e *GiteaExtractor) Extract(ctx context.Context, service *domain.Service) []domain.PullRequest {
	base := strings.TrimSuffix(service.URL, "/") + "/api/v1"

	issues, err := e.search(ctx, base, service)
	if err != nil {
		log.Print(err)
		return nil
	}

	var prs []domain.PullRequest
	for _, issue := range issues {
		repo := issue.Repository.FullName
		if len(service.Repositories) > 0 && !slices.Contains(service.Repositories, repo) {
			continue
		}

		var pull giteaPull
		if err := e.get(ctx, service, fmt.Sprintf("%s/repos/%s/pulls/%d", base, repo, issue.Number), &pull); err != nil {
			log.Print(err)
			continue
		}

		if len(service.Author) > 0 && !strings.EqualFold(pull.User.Login, service.Author) {
			continue
		}

		var status struct {
			State string `json:"state"`
		}
		if err := e.get(ctx, service, fmt.Sprintf("%s/repos/%s/commits/%s/status", base, repo, pull.Head.Sha), &status); err != nil {
			log.Print(err)
		}

		var reviews []giteaReview
		if err := e.get(ctx, service, fmt.Sprintf("%s/repos/%s/pulls/%d/reviews", base, repo, issue.Number), &reviews); err != nil {
			log.Print(err)
		}

		mergeable := "CONFLICTING"
		if pull.Mergeable {
			mergeable = "MERGEABLE"
		}

		prs = append(prs, domain.PullRequest{
			Repository:     pull.Base.Repo.FullName,
			RepositoryURL:  pull.Base.Repo.HTMLURL,
			Title:          pull.Title,
			Number:         pull.Number,
			Link:           pull.HTMLURL,
			CreatedAt:      pull.CreatedAt,
			UpdatedAt:      pull.UpdatedAt,
			Author:         pull.User.Login,
			IsDraft:        pull.Draft,
			ReviewDecision: giteaReviewDecision(reviews),
			Mergeable:      mergeable,
			CheckStatus:    strings.ToUpper(status.State),
		})
	}

	return prs
}

func (e *GiteaExtractor) search(ctx context.Context, base string, service *domain.Service) ([]giteaIssue, error) {
	query := url.Values{}
	query.Set("type", "pulls")
	query.Set("state", "open")
	query.Set("limit", fmt.Sprint(giteaPageSize))
	if len(service.Owner) > 0 {
		query.Set("owner", service.Owner)
	}

	var all []giteaIssue
	for page := 1; ; page++ {
		query.Set("page", fmt.Sprint(page))

		var issues []giteaIssue
		if err := e.get(ctx, service, base+"/repos/issues/search?"+query.Encode(), &issues); err != nil {
			return nil, err
		}

		all = append(all, issues...)
		if len(issues) < giteaPageSize {
			return all, nil
		}
	}
}

func (e *GiteaExtractor) get(ctx context.Context, service *domain.Service, endpoint string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if len(service.Token) > 0 {
		req.Header.Set("Authorization", "token "+service.Token)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gitea: GET %s: %s", endpoint, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func giteaReviewDecision(reviews []giteaReview) string {
	latest := make(map[string]string)
	for _, review := range reviews {
		if review.Dismissed {
			continue
		}
		if review.State == "APPROVED" || review.State == "REQUEST_CHANGES" {
			latest[review.User.Login] = review.State
		}
	}

	decision := ""
	for _, state := range latest {
		if state == "REQUEST_CHANGES" {
			return "CHANGES_REQUESTED"
		}
		decision = "APPROVED"
	}
	return decision
}

var _ Extractor = (*GiteaExtractor)(nil)
