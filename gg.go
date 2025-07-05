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

type CloneFlags struct {
	Owner           string
	Repo            string
	OutDir          string
	Shallow         bool
	RepoFile        string
	IncludeArchived bool
}

type KeyringProvider interface {
	Get() (string, error)
	Set(string) error
	Delete() error
}

type InputReader func() (string, error)

type GitHubService interface {
	Login(ctx context.Context) error
	Logout(ctx context.Context) error
	Show(ctx context.Context) error
	Clone(ctx context.Context, flags CloneFlags) error
}

type GitHubClient interface {
	GetAuthenticatedUser(ctx context.Context) (string, error)
	ListRepositoriesByUser(ctx context.Context, user string) ([]Repo, error)
	SearchRepositoriesByName(ctx context.Context, name string) ([]Repo, error)
	FindRepos(ctx context.Context, opts FindRepoOpts) ([]Repo, error)
}

type GitHubClientProvider func(token string) GitHubClient

type GitClient interface {
	Clone(ctx context.Context, repos []Repo, outDIr string, shallow bool) error
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
