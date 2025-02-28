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
Convenience cli for everyday things

Usage:
  gg [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  github      Convenience wrapper for github stuff
  help        Help about any command
  version     Print the version information

Flags:
      --debug   Enable debug logging
  -h, --help    help for gg

Use "gg [command] --help" for more information about a command.
```

## GitHub Authentication

GitHub Authentication is done by providing the CLI a GitHub token, which is stored in the system keyring and used for subsequent requests.

```bash
gg github login
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
