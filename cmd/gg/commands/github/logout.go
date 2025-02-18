package github

import (
	"fmt"

	"github.com/joakimen/gg/github"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Show stored GitHub credentials",
		RunE: func(_ *cobra.Command, _ []string) error {
			credentialsManager := github.NewCredentialsManager()
			err := credentialsManager.ClearToken()
			if err != nil {
				return fmt.Errorf("failed to clear existing github credentials from keyring: %w", err)
			}
			fmt.Println("GitHub credentials were cleared from the system keyring.")
			return nil
		},
	}
	return cmd
}
