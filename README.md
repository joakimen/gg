# gg

gg - Go GitHub

Convenience CLI for some GitHub operations.

## Requirements

- [gh](https://cli.github.com/)

## Installation

### Homebrew

```shell
brew install joakimen/tap/gg
```

### Go

```shell
go install github.com/joakimen/gg/cmd/gg@latest
```

## Usage

```shell
$ gg
Usage: gg <command> [flags]

Interactive GitHub repo cloning

Flags:
  -h, --help     Show context-sensitive help.
      --debug    Enable debug logging

Commands:
  version    Print version number
  clone      Clone one or more repos

Run "gg <command> --help" for more information on a command.

gg: error: expected one of "version", "clone"
```

## GitHub Authentication

This CLI currently uses the `gh` CLI to clone repositories. You will need to authenticate with GitHub before using this CLI.

Authenticating to GitHub:

```bash
gh auth login
```

## Development

### Git hooks

Git hooks are managed using [lefthook](https://github.com/evilmartians/lefthook).

1. Install lefthook: `brew install lefthook`
2. Run `lefthook install` in the repository root

At this point, hooks will run automatically on the Git hooks configured in `.lefthook.yml`, such as `pre-commit`.

To manually invoke the checks on all files without Git hooks, run:

```shell
$ lefthook run pre-commit --all-files
╭───────────────────────────────────────╮
│ 🥊 lefthook v1.10.1  hook: pre-commit │
╰───────────────────────────────────────╯
│  gofumpt -l -w . (skip) no matching staged files
│  goimports -w . (skip) no matching staged files
│  staticcheck ./... (skip) no matching staged files
│  go vet ./... (skip) no matching staged files
┃  yaml-lint ❯
┃  json-lint ❯
┃  md-lint ❯
  ────────────────────────────────────
summary: (done in 0.16 seconds)
✔️ yaml-lint
✔️ json-lint
✔️ md-lint
```
