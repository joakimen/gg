package github

import (
	"github.com/joakimen/gg/pkg/cmd/github/clone"
	"github.com/joakimen/gg/pkg/cmd/github/login"
	"github.com/joakimen/gg/pkg/cmd/github/logout"
	"github.com/joakimen/gg/pkg/cmd/github/show"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "Convenience wrapper for github stuff",
	}

	githubCmd.AddCommand(
		login.NewCmd(),
		show.NewCmd(),
		logout.NewCmd(),
		clone.NewCmd(),
	)

	return githubCmd
}
