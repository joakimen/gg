package github

import (
	"fmt"

	"github.com/joakimen/gg/github"
	"github.com/spf13/cobra"
)

func NewShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show stored GitHub credentials",
		RunE: func(_ *cobra.Command, _ []string) error {
			credentialsManager := github.NewCredentialsManager()
			token, err := credentialsManager.GetToken()
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
