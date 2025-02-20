package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/joakimen/gg"
)

type Service struct {
	Client *github.Client
}

func NewService(authToken string) Service {
	timeoutSeconds := 10
	httpClient := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	return Service{
		Client: github.NewClient(httpClient).WithAuthToken(authToken),
	}
}

func (s *Service) GetAuthenticatedUser(ctx context.Context) (gg.GitHubUser, error) {
	user, _, err := s.Client.Users.Get(ctx, "")
	if err != nil {
		return gg.GitHubUser{}, fmt.Errorf("failed to get the authenticated user: %w", err)
	}

	mappedUser := gg.GitHubUser{
		Login: user.GetLogin(),
	}
	return mappedUser, nil
}

func (s *Service) ListRepositoriesByUser(ctx context.Context, user string) ([]gg.Repo, error) {
	opts := &github.RepositoryListByUserOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allRepos []gg.Repo
	for {
		repos, resp, err := s.Client.Repositories.ListByUser(ctx, user, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories for user %s: %w", user, err)
		}
		for _, repo := range repos {
			repo := gg.Repo{
				Owner: repo.GetOwner().GetLogin(),
				Name:  repo.GetName(),
			}
			allRepos = append(allRepos, repo)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allRepos, nil
}

func (s *Service) SearchRepositoriesByName(ctx context.Context, name string) ([]gg.Repo, error) {
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allRepos []gg.Repo
	for {
		repos, resp, err := s.Client.Search.Repositories(ctx, name, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to search repositories by name %s: %w", name, err)
		}
		for _, repo := range repos.Repositories {
			repo := gg.Repo{
				Owner: repo.GetOwner().GetLogin(),
				Name:  repo.GetName(),
			}
			allRepos = append(allRepos, repo)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allRepos, nil
}
