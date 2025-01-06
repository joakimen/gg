package github

import (
	"fmt"
)

type Repo struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type RepoCloneResult struct {
	Repo Repo  `json:"repo"`
	Err  error `json:"err"`
}
type RepoSearchResult struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func (r Repo) NameWithOwner() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}
