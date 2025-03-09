package show

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
	token, err := c.Keyring.Get()
	if err != nil {
		return err
	}
	fmt.Println(token)
	return nil
}
