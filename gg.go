package gg

import "context"

type Repo struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CloneResult struct {
	Repo Repo  `json:"repo"`
	Err  error `json:"err"`
}

type KeyringService interface {
	Get() (string, error)
	Set(string) error
	Delete() error
}

type TTY interface {
	Read(string) (string, error)
}

type GitHubService interface {
	GetAuthenticatedUser(context.Context) (string, error)
	ListRepositoriesByUser(context.Context, string) ([]Repo, error)
	SearchRepositoriesByName(context.Context, string) ([]Repo, error)
}

type GitService interface {
	Clone(Repo, string, bool) error
}

type RepoFuzzyProvider interface {
	Select([]Repo) ([]Repo, error)
}
