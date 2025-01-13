package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	DefaultRepoListLimit = 100
)

type Config struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
	RepoFile        string
	InfoFn          func(...interface{})
	DebugFn         func(...interface{})
}

type Flags struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
	RepoFile        string
}

type Envs struct {
	CloneDir string
}

func Load(args []string, outWriter io.Writer) (Config, error) {
	// initialize cfg with defaults
	cfg := Config{
		Limit: DefaultRepoListLimit,
	}

	// parse flags
	flags := Flags{}
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	flag.StringVar(&flags.Owner, "o", cfg.Owner, "owner of the repository to clone")
	flag.StringVar(&flags.Repo, "r", cfg.Repo, "name of the repository to clone")
	flag.BoolVar(&flags.IncludeArchived, "a", cfg.IncludeArchived, "include archived repositories")
	flag.IntVar(&flags.Limit, "l", cfg.Limit, "limit the number of repositories to search for")
	flag.StringVar(&flags.CloneDir, "d", cfg.CloneDir, "directory to clone the repositories into")
	flag.StringVar(&flags.RepoFile, "f", cfg.RepoFile, "File containing the list of repositories to clone")
	flag.BoolVar(&flags.Verbose, "v", cfg.Verbose, "verbose output")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	// parse env vars
	envs := Envs{
		CloneDir: os.Getenv("CLONE_DIR"),
	}

	// conclude configurations
	cfg.Owner = flags.Owner
	cfg.Repo = flags.Repo
	cfg.IncludeArchived = flags.IncludeArchived
	cfg.Verbose = flags.Verbose
	cfg.RepoFile = flags.RepoFile
	cfg.Limit = flags.Limit

	switch {
	case flags.CloneDir != "":
		cfg.CloneDir = flags.CloneDir
	case envs.CloneDir != "":
		cfg.CloneDir = envs.CloneDir
	default:
		return Config{},
			errors.New("clone directory not specified, set either the CLONE_DIR environment variable or use the -d flag")
	}

	cfg.InfoFn = func(args ...interface{}) {
		fmt.Fprintln(outWriter, args...)
	}

	cfg.DebugFn = func(args ...interface{}) {
		if cfg.Verbose {
			fmt.Fprintln(outWriter, args...)
		}
	}

	return cfg, nil
}

func DirExists(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("dir doesn't exist: %s, %w", path, err)
	}
	return nil
}
