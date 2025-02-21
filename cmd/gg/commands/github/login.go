package github

import (
	"fmt"
	"log/slog"
	"syscall"

	"github.com/joakimen/gg/github"
	"github.com/joakimen/gg/keyring"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate to GitHub",
		RunE: func(cmd *cobra.Command, _ []string) error {
			slog.Debug("reading token from user")

			keyringManager := keyring.NewManager(keyringUser)
			inputToken, err := readPassword("Enter your GitHub API token: ")
			if err != nil {
				return fmt.Errorf("failed to read token: %w", err)
			}

			if inputToken == "" {
				return fmt.Errorf("token must not be empty")
			}

			fmt.Println("Verifying token..")
			githubService := github.NewService(inputToken)
			user, err := githubService.GetAuthenticatedUser(cmd.Context())
			if err != nil {
				return fmt.Errorf("authentication using the provided token failed: %w", err)
			}
			err = keyringManager.Set(inputToken)
			if err != nil {
				return fmt.Errorf("failed to store token in keyring: %w", err)
			}
			fmt.Println("Authenticated successfully as user:", user.Login)
			return nil
		},
	}
	return cmd
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytepw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()
	return string(bytepw), nil
}
