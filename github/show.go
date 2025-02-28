package github

import (
	"fmt"

	"github.com/joakimen/gg/keyring"
)

func Show() error {
	keyringManager := keyring.NewManager(KeyringUser)
	token, err := keyringManager.Get()
	if err != nil {
		fmt.Println("No existing credentials found in keyring.")
	} else {
		fmt.Println(token)
	}
	return nil
}
