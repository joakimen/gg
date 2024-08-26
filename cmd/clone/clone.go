package main

import (
	"flag"
	"fmt"
	"github.com/joakimen/clone/pkg/fuzzy"
	"github.com/joakimen/clone/pkg/github"
	"io"
	"log/slog"
	"os"
)

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

type Config struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
}

type Flags struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
}

type Envs struct {
	CloneDir string
}

func loadConfig() (Config, error) {
	cfg := Config{
		Limit: 100,
	}

	flags := Flags{}
	flag.StringVar(&flags.Owner, "o", cfg.Owner, "owner of the repository to clone")
	flag.StringVar(&flags.Repo, "r", cfg.Repo, "name of the repository to clone")
	flag.BoolVar(&flags.IncludeArchived, "a", cfg.IncludeArchived, "include archived repositories")
	flag.IntVar(&flags.Limit, "l", cfg.Limit, "limit the number of repositories to search for")
	flag.StringVar(&flags.CloneDir, "d", cfg.CloneDir, "directory to clone the repositories into")
	flag.BoolVar(&flags.Verbose, "v", cfg.Verbose, "verbose output")
	flag.Parse()

	envs := Envs{
		CloneDir: os.Getenv("CLONE_DIR"),
	}

	if flags.Owner != "" {
		cfg.Owner = flags.Owner
	}

	if flags.Repo != "" {
		cfg.Repo = flags.Repo
	}

	if flags.IncludeArchived {
		cfg.IncludeArchived = flags.IncludeArchived
	}

	if flags.Limit != 0 {
		cfg.Limit = flags.Limit
	}

	if envs.CloneDir != "" {
		cfg.CloneDir = envs.CloneDir
	} else if flags.CloneDir != "" {
		cfg.CloneDir = flags.CloneDir
	} else {
		return Config{}, fmt.Errorf("clone directory not specified, set either the CLONE_DIR environment variable or use the -d flag")
	}

	if flags.Verbose {
		cfg.Verbose = flags.Verbose
	}

	return cfg, nil
}

func configureLogger(w io.Writer, cfg Config) *slog.Logger {
	var handlerOpts slog.HandlerOptions
	if cfg.Verbose {
		handlerOpts.Level = slog.LevelInfo
	} else {
		handlerOpts.Level = slog.LevelWarn
	}
	logger := slog.New(slog.NewTextHandler(w, &handlerOpts))
	return logger
}

func run(w io.Writer, args []string) error {

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := configureLogger(w, cfg)
	logger.Info("done merging config", "config", cfg)

	info, err := os.Stat(cfg.CloneDir)
	if os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("the specified clone directory does not exist: %s", cfg.CloneDir)
	}

	var repos []github.Repo
	if cfg.Owner != "" && cfg.Repo != "" {
		logger.Debug("both owner and repo was specified, not searching for repos")
		repos = []github.Repo{
			{
				Owner: cfg.Owner,
				Name:  cfg.Repo,
			},
		}
	} else {
		logger.Info("searching for repos", "owner", cfg.Owner, "repo", cfg.Repo, "includeArchived", cfg.IncludeArchived, "limit", cfg.Limit)
		var err error
		repos, err = github.ListRepos(cfg.Owner, cfg.Repo, cfg.IncludeArchived, cfg.Limit)
		if err != nil {
			return err
		}
		logger.Info("repo search complete", "count", (len(repos)))
		logger.Debug("search returned repos", "repos", repos)
	}

	var selectedRepos []github.Repo
	if len(repos) == 1 {
		selectedRepos = repos
	} else {
		logger.Info("filtering repos")
		var err error
		selectedRepos, err = fuzzy.SelectRepos(repos)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(w, "%s\n", "cloning repos:")
	for _, repo := range selectedRepos {
		fmt.Fprintf(w, "- %s\n", repo.NameWithOwner())
	}

	type Result struct {
		repo github.Repo
		err  error
	}
	resultChan := make(chan Result, len(selectedRepos))
	for _, repo := range selectedRepos {
		go func(r github.Repo) {
			err := github.Clone(cfg.CloneDir, r)
			if err != nil {
				resultChan <- Result{repo: r, err: err}
				return
			}
			resultChan <- Result{repo: r, err: nil}
		}(repo)
	}

	// receive clone results from result channel
	var errorResults []Result
	for i := 0; i < len(selectedRepos); i++ {
		result := <-resultChan
		if result.err != nil {
			errorResults = append(errorResults, result)
		}
	}
	close(resultChan)

	fmt.Fprintf(w, "\n")
	if len(errorResults) > 0 {
		fmt.Fprintf(w, "some repos failed to clone:\n")
		for _, result := range errorResults {
			fmt.Fprintf(w, "- %s: %v\n", result.repo.NameWithOwner(), result.err)
		}
	} else {
		fmt.Fprintf(w, "all repos cloned successfully!\n")
	}
	return nil
}
