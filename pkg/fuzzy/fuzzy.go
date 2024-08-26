// Package fuzzy provides fuzzy selection functionality for github.Repo
package fuzzy

import (
	"github.com/joakimen/clone/pkg/github"
	fz "github.com/ktr0731/go-fuzzyfinder"
)

// SelectRepos Provides fuzzy multi-selection from a slice of [github.Repo]
func SelectRepos(repos []github.Repo) ([]github.Repo, error) {
	renderFunc := func(selectedIndex int) string {
		return repos[selectedIndex].NameWithOwner()
	}
	indices, err := fz.FindMulti(repos, renderFunc)
	if err != nil {
		return nil, err
	}

	var selectedRepos []github.Repo
	for _, idx := range indices {
		selectedRepos = append(selectedRepos, repos[idx])
	}
	return selectedRepos, nil
}
