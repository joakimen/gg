package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joakimen/gg"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Clone(repo gg.Repo, outDir string, shallow bool) error {
	repoAbsPath := filepath.Join(outDir, repo.Owner, repo.Name)
	if _, err := os.Stat(repoAbsPath); !os.IsNotExist(err) {
		return fmt.Errorf("repo %s already exists in %s", repo.NameWithOwner(), repoAbsPath)
	}

	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", repo.Owner, repo.Name)

	args := []string{"clone", cloneURL, outDir}

	if shallow {
		args = append(args, "--depth", "1")
	}

	cloneCmd := exec.Command("git", args...)
	return cloneCmd.Run()
}
