package main

import "log/slog"

type VersionCommand struct{}

func (c *VersionCommand) Run(cfg Config) error {

	slog.Debug(version)
	return nil
}
