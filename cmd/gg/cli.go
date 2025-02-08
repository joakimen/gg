package main

const (
	cliName = "gg"
	cliDesc = "Convenience CLI for some GitHub operations"
)

// overwritten by goreleaser.
var version = "(development build)"

type CLI struct {
	Globals
	Version VersionCmd `cmd:"" help:"Print version number"`
	Clone   CloneCmd   `cmd:"" help:"Clone one or more repos"`
}

type Globals struct {
	Debug bool `help:"Enable debug logging"`
}

type VersionCmd struct{}

type CloneCmd struct {
	Repo            string `short:"r" help:"Name of the repository to clone"`
	Owner           string `short:"o" help:"Owner of the repository to clone"`
	IncludeArchived bool   `short:"a" help:"Include archived repositories"`
	Limit           int    `short:"l" default:"100" help:"Limit the number of repositories to search for"`
	CloneDir        string `short:"d" required:"true" env:"GG_CLONE_DIR" help:"Directory to clone the repositories into"`
	RepoFile        string `short:"f" help:"File containing the list of repositories to clone"`
	Shallow         bool   `help:"Shallow clone the repository"`
}
