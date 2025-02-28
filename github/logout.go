package github

import (
	"fmt"

	"github.com/joakimen/gg/keyring"
)

func Logout() error {
	keyringManager := keyring.NewManager(KeyringUser)
	err := keyringManager.Delete()
	if err != nil {
		return fmt.Errorf("failed to clear existing github credentials from keyring: %w", err)
	}
	fmt.Println("GitHub credentials were cleared from the system keyring.")
	return nil
}
