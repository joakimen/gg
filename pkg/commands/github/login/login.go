package login

import (
	"context"
	"fmt"
	"strings"

	"github.com/joakimen/gg"
)

var _ gg.GitHubCommand = (*Command)(nil)

type Command struct {
	Keyring       gg.KeyringProvider
	TTY           gg.TTYProvider
	ClientFactory gg.GitHubClientFactory
}

func (c *Command) Run(ctx context.Context) error {
	// read api token from user
	token, err := c.TTY.Read("Enter your GitHub API token: ")
	if err != nil {
		return err
	}

	if strings.TrimSpace(token) == "" {
		return fmt.Errorf("the provided token cannot be empty")
	}

	// test user token
	client := c.ClientFactory(token)
	userLogin, err := client.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Successfully authenticated as user:", userLogin)

	// if token is valid, store it
	err = c.Keyring.Set(token)
	if err != nil {
		return err
	}
	fmt.Println("Token successfully stored in system keyring.")

	return nil
}
