package fuzzy

import (
	"fmt"

	"github.com/joakimen/gg"
	fz "github.com/ktr0731/go-fuzzyfinder"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// Select provides fuzzy single- or multi-selection for repos.
func (s *Service) Select(repos []gg.Repo) ([]gg.Repo, error) {
	renderFunc := func(selectedIndex int) string {
		return repos[selectedIndex].NameWithOwner()
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
