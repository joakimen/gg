package git

import (
	"fmt"
	"log/slog"
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
	outDirAbs := filepath.Join(outDir, repo.Owner, repo.Name)
	if _, err := os.Stat(outDirAbs); !os.IsNotExist(err) {
		return fmt.Errorf("repo %s already exists in %s", repo.NameWithOwner(), outDirAbs)
	}

	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", repo.Owner, repo.Name)
	args := []string{"clone", cloneURL, outDirAbs}

	if shallow {
		args = append(args, "--depth", "1")
	}

	slog.Debug("cloning repo", "args", args)
	cloneCmd := exec.Command("git", args...)
	return cloneCmd.Run()
}
