package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v69/github"
)

const (
	repoPerPage          = 100
	clientTimeoutSeconds = 10
)

type Client struct {
	api *github.Client
}

func NewClient(token string) *Client {
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

// func (c *Client) ListRepositoriesByUser(ctx context.Context, user string) ([]gg.Repo, error) {
// 	opts := &github.RepositoryListByUserOptions{
// 		ListOptions: github.ListOptions{PerPage: repoPerPage},
// 	}
// 	var allRepos []gg.Repo
// 	for {
// 		repos, resp, err := c.api.Repositories.ListByUser(ctx, user, opts)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to list repositories for user %s: %w", user, err)
// 		}
// 		for _, repo := range repos {
// 			repo := gg.Repo{
// 				Owner: repo.GetOwner().GetLogin(),
// 				Name:  repo.GetName(),
// 			}
// 			allRepos = append(allRepos, repo)
// 		}
// 		if resp.NextPage == 0 {
// 			break
// 		}
// 		opts.Page = resp.NextPage
// 	}
// 	return allRepos, nil
// }

// func (c *Client) SearchRepositoriesByName(ctx context.Context, name string) ([]gg.Repo, error) {
// 	opts := &github.SearchOptions{
// 		ListOptions: github.ListOptions{PerPage: repoPerPage},
// 	}
// 	var allRepos []gg.Repo
// 	for {
// 		repos, resp, err := c.api.Search.Repositories(ctx, name, opts)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to search repositories by name %s: %w", name, err)
// 		}
// 		for _, repo := range repos.Repositories {
// 			repo := gg.Repo{
// 				Owner: repo.GetOwner().GetLogin(),
// 				Name:  repo.GetName(),
// 			}
// 			allRepos = append(allRepos, repo)
// 		}
// 		if resp.NextPage == 0 {
// 			break
// 		}
// 		opts.Page = resp.NextPage
// 	}
// 	return allRepos, nil
// }

// func (c *Client) Clone(
// 	ctx context.Context,
// 	git gg.GitService,
// 	repos []gg.Repo,
// 	outDir string,
// 	shallow bool,
// ) error {
// 	slog.DebugContext(ctx, "cloning repositories")
// 	if outDir == "" {
// 		return fmt.Errorf("refusing to clone without a specified outDir")
// 	}

// 	if len(repos) == 0 {
// 		return fmt.Errorf("no repos to clone")
// 	}

// 	fmt.Printf("cloning %d repos(s)..\n", len(repos))
// 	slog.DebugContext(ctx, "cloning repos", "outDir", outDir, "repos", repos)

// 	var wg sync.WaitGroup
// 	resultChan := make(chan gg.CloneResult, len(repos))
// 	for _, repo := range repos {
// 		wg.Add(1)
// 		go func(r gg.Repo) {
// 			defer wg.Done()
// 			cloneError := git.Clone(repo, outDir, shallow)
// 			resultChan <- gg.CloneResult{Repo: r, Err: cloneError}
// 		}(repo)
// 	}
// 	wg.Wait()
// 	close(resultChan)
// 	var errs []gg.CloneResult
// 	for res := range resultChan {
// 		if res.Err != nil {
// 			errs = append(errs, res)
// 		}
// 	}

// 	if len(errs) > 0 {
// 		fmt.Fprintln(os.Stderr, "failed to clone some repos:")
// 		for _, e := range errs {
// 			fmt.Fprintf(os.Stderr, "- %s/%s: %v\n", e.Repo.Owner, e.Repo.Name, e.Err)
// 		}
// 	} else {
// 		fmt.Println("all repos cloned successfully!")
// 	}
// 	return nil
// }

// type FindRepoOpts struct {
// 	RepoFilter        gg.RepoFuzzyProvider
// 	Owner             string
// 	Repo              string
// 	OutDir            string
// 	Shallow           bool
// 	RepoFile          string
// 	DefaultGitHubUser string
// }

// func (c *Client) FindRepos(ctx context.Context, opts FindRepoOpts) ([]gg.Repo, error) {
// 	// If a repo file is provided, unmarshal its content.
// 	if opts.RepoFile != "" {
// 		repoJSONData, err := os.ReadFile(opts.RepoFile)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to read repo file: %w", err)
// 		}
// 		var repos []gg.Repo
// 		if err = json.Unmarshal(repoJSONData, &repos); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal repos from file: %w", err)
// 		}
// 		return repos, nil
// 	}

// 	// If both owner and repo are provided, clone that repo.
// 	if opts.Owner != "" && opts.Repo != "" {
// 		return []gg.Repo{{Owner: opts.Owner, Name: opts.Repo}}, nil
// 	}

// 	// If only owner is provided, list repos and let user select.
// 	if opts.Owner != "" {
// 		slog.DebugContext(ctx, "listing repos by owner", "owner", opts.Owner)
// 		allRepos, err := c.ListRepositoriesByUser(ctx, opts.Owner)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to list repos by owner: %w", err)
// 		}
// 		repos, err := opts.RepoFilter.Select(allRepos)
// 		if err != nil {
// 			return nil, fmt.Errorf("error while filtering repos: %w", err)
// 		}
// 		return repos, nil
// 	}

// 	// If only repo is provided, search and let user select.
// 	if opts.Repo != "" {
// 		slog.DebugContext(ctx, "searching repos by name", "name", opts.Repo)
// 		allRepos, err := c.SearchRepositoriesByName(ctx, opts.Repo)
// 		if err != nil {
// 			return nil, err
// 		}
// 		repos, err := opts.RepoFilter.Select(allRepos)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return repos, nil
// 	}

// 	// Default: use the default GitHub user from env.
// 	if opts.DefaultGitHubUser == "" {
// 		return nil, fmt.Errorf("no default GitHub user configured")
// 	}

// 	allRepos, err := c.ListRepositoriesByUser(ctx, opts.DefaultGitHubUser)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return opts.RepoFilter.Select(allRepos)
// }
