# clone

Clones one or more repositories to the configured clone directory.

## Requirements

- [gh](https://cli.github.com/)

## Usage

```shell
$ clone 

    # Optional flags:
    -o  The owner of the repositories to clone
    -r  The name repositories to clone
    -d  The directory to clone the repositories to
    -a  Include archived repositories in the search
    -l  The limit of repositories to search for
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

This script uses the `gh` CLI to clone repositories. You will need to authenticate with GitHub before using this script.

Authenticating to GitHub:

```bash
gh auth login
```

## Ideas for later

- Add support for a config file that lets us specify a list of repositories to clone that we can share across machines