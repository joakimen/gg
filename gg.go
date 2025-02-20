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

type GitHubUser struct {
	Login string
}
