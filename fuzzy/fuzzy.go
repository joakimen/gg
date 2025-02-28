package fuzzy

import (
	"fmt"

	"github.com/joakimen/gg"
	fz "github.com/ktr0731/go-fuzzyfinder"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

// Select provides fuzzy single- or multi-selection for repos.
func (s *Provider) Select(repos []gg.Repo) ([]gg.Repo, error) {
	renderFunc := func(selectedIndex int) string {
		repo := repos[selectedIndex]
		return fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
	}
	indices, err := fz.FindMulti(repos, renderFunc)
	if err != nil {
		return nil, fmt.Errorf("no repos selected: %w", err)
	}

	var selectedRepos []gg.Repo
	for _, idx := range indices {
		selectedRepos = append(selectedRepos, repos[idx])
	}
	return selectedRepos, nil
}
