package fuzzy

import (
	"fmt"
	"github.com/joakimen/clone/pkg/github"
	fz "github.com/ktr0731/go-fuzzyfinder"
)

// SelectRepos Provides fuzzy multi-selection from a slice of [github.Repo]
func SelectRepos(repos []github.Repo) ([]github.Repo, error) {
	renderFunc := func(selectedIndex int) string {
		return repos[selectedIndex].NameWithOwner()
	}

	previewFunc := func(selectedIndex, width, height int) string {
		if selectedIndex == -1 {
			return ""
		}
		return fmt.Sprintf("Cur repo: %s\n", repos[selectedIndex].NameWithOwner())
	}

	indices, err := fz.FindMulti(repos, renderFunc, fz.WithPreviewWindow(previewFunc))
	if err != nil {
		return nil, err
	}

	var selectedRepos []github.Repo
	for _, idx := range indices {
		selectedRepos = append(selectedRepos, repos[idx])
	}
	return selectedRepos, nil
}
