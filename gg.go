package gg

import (
	"context"
)

type Repo struct {
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	Archived bool   `json:"archived"`
}

type CloneResult struct {
	Repo Repo  `json:"repo"`
	Err  error `json:"err"`
}

type KeyringProvider interface {
	Get() (string, error)
	Set(string) error
	Delete() error
}

type TTYProvider interface {
	Read(string) (string, error)
}

type GitHubClient interface {
	GetAuthenticatedUser(ctx context.Context) (string, error)
	ListRepositoriesByUser(ctx context.Context, user string) ([]Repo, error)
	SearchRepositoriesByName(ctx context.Context, name string) ([]Repo, error)
	Clone(ctx context.Context, git GitClient, repos []Repo, outDir string, shallow bool) error
	FindRepos(ctx context.Context, opts FindRepoOpts) ([]Repo, error)
}

type GitHubClientProvider func(token string) GitHubClient

type GitClient interface {
	Clone(Repo, string, bool) error
}

type RepoSelector interface {
	Select([]Repo) ([]Repo, error)
}

type FindRepoOpts struct {
	RepoSelector      RepoSelector
	Owner             string
	Repo              string
	OutDir            string
	Shallow           bool
	RepoFile          string
	DefaultGitHubUser string
	IncludeArchived   bool
}
