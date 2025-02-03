package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	DefaultRepoListLimit = 100
)

type Config struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	RepoListLimit   int
	CloneDir        string
	Verbose         bool
	RepoFile        string
	DebugLogging    bool
	ShallowClone    bool
}

type flags struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	RepoListLimit   int
	CloneDir        string
	Verbose         bool
	RepoFile        string
	DebugLogging    bool
	ShallowClone    bool
}

type envs struct {
	CloneDir string
}

func newConfig() Config {
	return Config{
		Owner:           "",
		Repo:            "",
		IncludeArchived: false,
		RepoListLimit:   DefaultRepoListLimit,
		CloneDir:        "",
		Verbose:         false,
		RepoFile:        "",
		DebugLogging:    false,
		ShallowClone:    false,
	}
}

func parseConfig(flags flags, envs envs) Config {

	cfg := newConfig()

	if flags.Owner != "" {
		cfg.Owner = flags.Owner
	}

	if flags.Repo != "" {
		cfg.Repo = flags.Repo
	}

	if flags.IncludeArchived {
		cfg.IncludeArchived = flags.IncludeArchived
	}

	if flags.IncludeArchived {
		cfg.IncludeArchived = flags.IncludeArchived
	}

	if flags.Verbose {
		cfg.Verbose = flags.Verbose
	}

	if flags.RepoFile != "" {
		cfg.RepoFile = flags.RepoFile
	}

	if flags.RepoListLimit != 0 {
		cfg.RepoListLimit = flags.RepoListLimit
	}

	if flags.ShallowClone {
		cfg.ShallowClone = flags.ShallowClone
	}

	if flags.DebugLogging {
		cfg.DebugLogging = flags.DebugLogging
	}

	switch {
	case flags.CloneDir != "":
		cfg.CloneDir = flags.CloneDir
	case envs.CloneDir != "":
		cfg.CloneDir = envs.CloneDir
	}

	return cfg

}

func parseFlags(args []string) (flags, error) {

	f := flags{}
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	fs.StringVar(&f.Owner, "o", "", "owner of the repository to clone")
	fs.StringVar(&f.Repo, "r", "", "name of the repository to clone")
	fs.BoolVar(&f.IncludeArchived, "a", false, "include archived repositories")
	fs.IntVar(&f.RepoListLimit, "l", DefaultRepoListLimit, "limit the number of repositories to search for")
	fs.StringVar(&f.CloneDir, "d", "", "directory to clone the repositories into")
	fs.StringVar(&f.RepoFile, "f", "", "File containing the list of repositories to clone")
	fs.BoolVar(&f.Verbose, "v", false, "verbose output")
	fs.BoolVar(&f.ShallowClone, "shallow", false, "shallow clone the repository")
	fs.BoolVar(&f.DebugLogging, "debug", false, "enable debug logging")

	if err := fs.Parse(args); err != nil {
		return flags{}, err
	}
	return f, nil
}

func LoadConfig(args []string) (Config, error) {

	// parse flags
	flags, err := parseFlags(args)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing flags: %w", err)
	}

	// parse env vars
	envs := envs{
		CloneDir: os.Getenv("CLONE_DIR"),
	}

	// parse final config
	cfg := parseConfig(flags, envs)

	return cfg, nil
}
