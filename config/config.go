package config

import (
	"flag"
	"fmt"
	"github.com/joakimen/clone"
	"os"
)

func Load() (clone.Config, error) {
	cfg := clone.Config{
		Limit: 100,
	}

	flags := clone.Flags{}
	flag.StringVar(&flags.Owner, "o", cfg.Owner, "owner of the repository to clone")
	flag.StringVar(&flags.Repo, "r", cfg.Repo, "name of the repository to clone")
	flag.BoolVar(&flags.IncludeArchived, "a", cfg.IncludeArchived, "include archived repositories")
	flag.IntVar(&flags.Limit, "l", cfg.Limit, "limit the number of repositories to search for")
	flag.StringVar(&flags.CloneDir, "d", cfg.CloneDir, "directory to clone the repositories into")
	flag.StringVar(&flags.RepoFile, "f", cfg.RepoFile, "File containing the list of repositories to clone")
	flag.BoolVar(&flags.Verbose, "v", cfg.Verbose, "verbose output")
	flag.Parse()

	envs := clone.Envs{
		CloneDir: os.Getenv("CLONE_DIR"),
	}

	cfg.Owner = flags.Owner
	cfg.Repo = flags.Repo
	cfg.IncludeArchived = flags.IncludeArchived
	cfg.Verbose = flags.Verbose
	cfg.RepoFile = flags.RepoFile
	cfg.Limit = flags.Limit

	if flags.CloneDir != "" {
		cfg.CloneDir = flags.CloneDir
	} else if envs.CloneDir != "" {
		cfg.CloneDir = envs.CloneDir
	} else {
		return clone.Config{}, fmt.Errorf("clone directory not specified, set either the CLONE_DIR environment variable or use the -d flag")
	}

	return cfg, nil
}
