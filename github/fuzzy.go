package github

import (
	"fmt"

	fz "github.com/ktr0731/go-fuzzyfinder"
)

// FuzzySelect provides fuzzy single- or multi-selection for repos.
func FuzzySelect(repos []Repo) ([]Repo, error) {
	renderFunc := func(selectedIndex int) string {
		return repos[selectedIndex].NameWithOwner()
	}
	indices, err := fz.FindMulti(repos, renderFunc)
	if err != nil {
		return nil, fmt.Errorf("no repos selected: %w", err)
	}

	var selectedRepos []Repo
	for _, idx := range indices {
		selectedRepos = append(selectedRepos, repos[idx])
	}
	return selectedRepos, nil
}
