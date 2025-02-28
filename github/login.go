package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/joakimen/gg/keyring"
	"github.com/joakimen/gg/tty"
)

func Login(ctx context.Context) error {
	slog.DebugContext(ctx, "reading token from user")
	keyringManager := keyring.NewManager(KeyringUser)

	passwordProvider := tty.NewPasswordProvider()
	inputToken, err := passwordProvider.ReadPassword("Enter your GitHub API token: ")
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}

	if inputToken == "" {
		return fmt.Errorf("token must not be empty")
	}

	fmt.Println("Verifying token..")
	githubClient := NewClient(inputToken)
	user, err := githubClient.GetAuthenticatedUser(ctx)
	if err != nil {
		return fmt.Errorf("authentication using the provided token failed: %w", err)
	}
	err = keyringManager.Set(inputToken)
	if err != nil {
		return fmt.Errorf("failed to store token in keyring: %w", err)
	}
	fmt.Println("Authenticated successfully as user:", user)
	return nil
}
