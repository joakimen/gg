package clone

import "fmt"

type Repo struct {
	Owner string
	Name  string
}

type Repos []Repo

type RepoCloneResult struct {
	Repo Repo
	Err  error
}
type RepoCloneResults []RepoCloneResult

type RepoSearchResult struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type RepoSearchResults []RepoSearchResult

func (r Repo) NameWithOwner() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}
