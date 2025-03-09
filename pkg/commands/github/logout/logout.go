package logout

import (
	"context"
	"fmt"

	"github.com/joakimen/gg"
)

var _ gg.GitHubCommand = (*Command)(nil)

type Command struct {
	Keyring gg.KeyringProvider
}

func (c *Command) Run(_ context.Context) error {
	fmt.Println("Deleting stored github credentials from system keyring..")
	err := c.Keyring.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete github credentals: %w", err)
	}
	fmt.Println("Done.")
	return nil
}
