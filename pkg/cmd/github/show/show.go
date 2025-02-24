package show

import (
	"fmt"

	"github.com/joakimen/gg/github"
	"github.com/joakimen/gg/keyring"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show stored GitHub credentials",
		RunE: func(_ *cobra.Command, _ []string) error {
			keyringManager := keyring.NewManager(github.KeyringUser)
			token, err := keyringManager.Get()
			if err != nil {
				fmt.Println("No existing credentials found in keyring.")
			} else {
				fmt.Println(token)
			}
			return nil
		},
	}
	return cmd
}
