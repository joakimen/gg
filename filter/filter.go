// Package filter provides functionality for filtering a list of [clone.Repo] object using fuzzy selection
package filter

import (
	"github.com/joakimen/clone"
	fz "github.com/ktr0731/go-fuzzyfinder"
)

// Select Provides fuzzy multi-selection from a slice of [clone.Repo]
func Select(repos clone.Repos) ([]clone.Repo, error) {
	renderFunc := func(selectedIndex int) string {
		return repos[selectedIndex].NameWithOwner()
	}
	indices, err := fz.FindMulti(repos, renderFunc)
	if err != nil {
		return nil, err
	}

	var selectedRepos clone.Repos
	for _, idx := range indices {
		selectedRepos = append(selectedRepos, repos[idx])
	}
	return selectedRepos, nil
}
