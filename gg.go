package gg

import (
	"fmt"
)

type Repo struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CloneResult struct {
	Repo Repo  `json:"repo"`
	Err  error `json:"err"`
}

func (r Repo) NameWithOwner() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

// ype GitHubCredentials struct {
// 	Username string `json:"username"`
// 	APIToken string `json:"api_token"`
// 	Host     string `json:"host"`
// }

type KeyringItem interface {
	Get() (string, error)
	Set(string) error
	Delete() error
}

type GitHubService interface {
	GetRepos(owner string, repo string, includeArchived bool, limit int) ([]Repo, error)
}

type GitHubUser struct {
	Login string
}
