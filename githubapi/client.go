package githubapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/joakimen/gg"
)

const (
	repoPerPage          = 100
	clientTimeoutSeconds = 10
)

// ensure types satisfy their interfaces.
var (
	_ gg.GitHubClient = (*Client)(nil)
)

type Client struct {
	api *github.Client
}

// TokenClientProvider is a function used to provide instantiation of a GitHub client outside of
// this package to preserve isolation of the go-github dependency in this package.
var TokenClientProvider gg.GitHubClientProvider = func(token string) gg.GitHubClient {
	httpClient := &http.Client{
		Timeout: time.Duration(clientTimeoutSeconds) * time.Second,
	}
	githubClient := github.NewClient(httpClient).WithAuthToken(token)
	return &Client{
		api: githubClient,
	}
}

func (c *Client) GetAuthenticatedUser(ctx context.Context) (string, error) {
	user, _, err := c.api.Users.Get(ctx, "")
	if err != nil {
		return "", fmt.Errorf("failed to get the authenticated user: %w", err)
	}
	return user.GetLogin(), nil
}

func (c *Client) ListRepositoriesByUser(ctx context.Context, user string) ([]gg.Repo, error) {
	opts := &github.RepositoryListByUserOptions{
		ListOptions: github.ListOptions{PerPage: repoPerPage},
	}
	slog.DebugContext(ctx, "Listing repositories for user", "user", user, "opts", opts)
	var allRepos []gg.Repo
	for {
		repos, resp, err := c.api.Repositories.ListByUser(ctx, user, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories for user %s: %w", user, err)
		}
		for _, repo := range repos {
			repo := gg.Repo{
				Owner:    repo.GetOwner().GetLogin(),
				Name:     repo.GetName(),
				Archived: repo.GetArchived(),
			}
			allRepos = append(allRepos, repo)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	slog.DebugContext(ctx, "returning found repos", "repo_count", len(allRepos))
	return allRepos, nil
}

func (c *Client) SearchRepositoriesByName(ctx context.Context, name string) ([]gg.Repo, error) {
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: repoPerPage},
	}
	var allRepos []gg.Repo
	for {
		repos, resp, err := c.api.Search.Repositories(ctx, name, opts)
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

func (c *Client) FindRepos(ctx context.Context, opts gg.FindRepoOpts) ([]gg.Repo, error) {
	// If a repo file is provided, unmarshal its content.
	if opts.RepoFile != "" {
		repoJSONData, err := os.ReadFile(opts.RepoFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read repo file: %w", err)
		}
		var repos []gg.Repo
		if err = json.Unmarshal(repoJSONData, &repos); err != nil {
			return nil, fmt.Errorf("failed to unmarshal repos from file: %w", err)
		}
		return repos, nil
	}

	// If both owner and repo are provided, clone that repo.
	if opts.Owner != "" && opts.Repo != "" {
		return []gg.Repo{{Owner: opts.Owner, Name: opts.Repo}}, nil
	}

	// If only owner is provided, list repos and let user select.
	if opts.Owner != "" {
		slog.DebugContext(ctx, "listing repos by owner", "owner", opts.Owner)
		allRepos, err := c.ListRepositoriesByUser(ctx, opts.Owner)
		if err != nil {
			return nil, fmt.Errorf("failed to list repos by owner: %w", err)
		}

		var includedRepos []gg.Repo
		for _, repo := range allRepos {
			if opts.IncludeArchived || !repo.Archived {
				includedRepos = append(includedRepos, repo)
			}
		}

		repos, err := opts.RepoSelector.Select(includedRepos)
		if err != nil {
			return nil, fmt.Errorf("error while filtering repos: %w", err)
		}
		return repos, nil
	}

	// If only repo is provided, search and let user select.
	if opts.Repo != "" {
		slog.DebugContext(ctx, "searching repos by name", "name", opts.Repo)
		allRepos, err := c.SearchRepositoriesByName(ctx, opts.Repo)
		if err != nil {
			return nil, err
		}
		repos, err := opts.RepoSelector.Select(allRepos)
		if err != nil {
			return nil, err
		}
		return repos, nil
	}

	// Default: use the default GitHub user from env.
	if opts.DefaultGitHubUser == "" {
		return nil, fmt.Errorf("no default GitHub user configured")
	}

	allRepos, err := c.ListRepositoriesByUser(ctx, opts.DefaultGitHubUser)
	if err != nil {
		return nil, err
	}

	return opts.RepoSelector.Select(allRepos)
}
