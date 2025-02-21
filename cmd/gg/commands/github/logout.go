package github

import (
	"fmt"

	"github.com/joakimen/gg/keyring"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Clear stored GitHub credentials",
		RunE: func(_ *cobra.Command, _ []string) error {
			keyringManager := keyring.NewManager(keyringUser)
			err := keyringManager.Delete()
			if err != nil {
				return fmt.Errorf("failed to clear existing github credentials from keyring: %w", err)
			}
			fmt.Println("GitHub credentials were cleared from the system keyring.")
			return nil
		},
	}
	return cmd
}
