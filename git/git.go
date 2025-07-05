package git

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/joakimen/gg"
)

var _ gg.GitClient = (*Client)(nil)

// single function interface?

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

type CloneOpts struct {
	Shallow bool
	OutDir  string
}

func (c *Client) Clone(
	ctx context.Context,
	repos []gg.Repo,
	outDir string,
	shallow bool,
) error {
	slog.DebugContext(ctx, "cloning repositories")
	if outDir == "" {
		return fmt.Errorf("refusing to clone without a specified outDir")
	}

	if len(repos) == 0 {
		return fmt.Errorf("no repos to clone")
	}

	fmt.Printf("cloning %d repos(s)..\n", len(repos))
	slog.DebugContext(ctx, "cloning repos", "outDir", outDir, "repos", repos)

	var wg sync.WaitGroup
	resultChan := make(chan gg.CloneResult, len(repos))
	for _, repo := range repos {
		wg.Add(1)
		go func(r gg.Repo) {
			defer wg.Done()
			cloneError := clone(ctx, repo, outDir, shallow)
			resultChan <- gg.CloneResult{Repo: r, Err: cloneError}
		}(repo)
	}
	wg.Wait()
	close(resultChan)
	var errs []gg.CloneResult
	for res := range resultChan {
		if res.Err != nil {
			errs = append(errs, res)
		}
	}

	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, "failed to clone some repos:")
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "- %s/%s: %v\n", e.Repo.Owner, e.Repo.Name, e.Err)
		}
	} else {
		fmt.Println("all repos cloned successfully!")
	}
	return nil
}

func clone(_ context.Context, repo gg.Repo, outDir string, shallow bool) error {
	outDirAbs := filepath.Join(outDir, repo.Owner, repo.Name)
	if _, err := os.Stat(outDirAbs); !os.IsNotExist(err) {
		return fmt.Errorf("repo %s/%s already exists in %s", repo.Owner, repo.Name, outDirAbs)
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
