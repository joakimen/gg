package main

type VersionCommand struct{}

func (c *VersionCommand) Run(cfg Config) error {
	info := cfg.InfoFn
	info(version)
	return nil
}
