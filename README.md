# clone

Clones one or more repositories to the configured clone directory.

## Requirements

- [gh](https://cli.github.com/)

## Usage

```shell
$ clone 

    # Optional flags:
    -o  Owner of the repositories to clone
    -r  Name repository to clone
    -d  Directory to clone the repositories to
    -a  Include archived repositories in the search
    -l  The limit of repositories to search for
    -f  Name of a file containing a list of repositories to clone
    -v  Verbose output
```

**Note**: All flags are optional, but the Clone Directory must be set either as an environment variable or as an argument:

```shell
# using env 
export CLONE_DIR=~/dev/github.com # (put this in your shell profile)
clone

# using arg
clone -d ~/dev/github.com
```

## GitHub Authentication

This CLI uses the `gh` CLI to clone repositories. You will need to authenticate with GitHub before using this CLI.

Authenticating to GitHub:

```bash
gh auth login
```
