package processor

import (
	"context"
	"fmt"
	"github.com/rapatao/pr-checker-go/domain"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"log"
	"time"
)

func extractGitHub(ctx context.Context, service *domain.Service) []domain.PullRequest {
	token := &oauth2.Token{AccessToken: service.Token}
	tokenSource := oauth2.StaticTokenSource(token)
	oauth2Client := oauth2.NewClient(ctx, tokenSource)
	client := githubv4.NewClient(oauth2Client)

	var prs []domain.PullRequest

	var query struct {
		Search struct {
			Edges []struct {
				Node struct {
					PullRequest struct {
						Repository struct {
							Url           string
							NameWithOwner string
						}

						Author struct {
							Login string
						}

						State     string
						CreatedAt time.Time
						UpdatedAt time.Time
						Number    int
						Title     string
						IsDraft   bool
						Url       string
					} `graphql:"... on PullRequest"`
				}
			}
		} `graphql:"search(query:$filter, type: ISSUE, last: 100)"`
	}

	filter := "type:pr state:open"

	if len(service.Author) > 0 {
		filter = fmt.Sprintf("%s author:%s", filter, service.Author)
	}

	for _, repository := range service.Repositories {
		filter = fmt.Sprintf("%s repo:%s", filter, repository)
	}

	variables := map[string]interface{}{
		"filter": githubv4.String(filter),
	}
	err := client.Query(ctx, &query, variables)
	if err != nil {
		log.Print(err)
	}

	for _, edge := range query.Search.Edges {
		pr := edge.Node.PullRequest

		if pr.Title == "" ||
			pr.Repository.NameWithOwner == "" ||
			pr.Url == "" {
			continue
		}

		prs = append(prs, domain.PullRequest{
			Service:       service.Name,
			Repository:    pr.Repository.NameWithOwner,
			RepositoryURL: pr.Repository.Url,
			Title:         pr.Title,
			Number:        pr.Number,
			Link:          pr.Url,
			CreatedAt:     pr.CreatedAt,
			UpdatedAt:     pr.UpdatedAt,
			Author:        pr.Author.Login,
			IsDraft:       pr.IsDraft,
		})

	}

	return prs
}
