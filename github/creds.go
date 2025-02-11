package github

import (
	"encoding/json"
	"fmt"

	"github.com/joakimen/gg"
)

const keyringItem = "github"

// CredentialsService manages GitHub credentials in the system keyring
type CredentialsService struct {
	keyring gg.KeyringItemProvider
}

type Credentials struct {
	Username string
	APIToken string
	Host     string
}

func NewAuthService(keychain gg.KeyringItemProvider) *CredentialsService {
	return &CredentialsService{
		keyring: keychain,
	}
}

func (cs *CredentialsService) Set(creds Credentials) error {
	credsJson, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to encode credentials: %w", err)
	}
	err = cs.keyring.Set(keyringItem, string(credsJson))
	if err != nil {
		return fmt.Errorf("failed to save credentials to keyring: %w", err)
	}
	return nil
}

func (cs *CredentialsService) Get() (Credentials, error) {
	credsJson, err := cs.keyring.Get(keyringItem)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to get credentials from keyring: %w", err)
	}
	var creds Credentials
	err = json.Unmarshal([]byte(credsJson), &creds)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to decode credentials: %w", err)
	}
	return creds, nil
}

func (cs *CredentialsService) Delete() error {
	err := cs.keyring.Delete(keyringItem)
	if err != nil {
		return fmt.Errorf("failed to delete credentials from keyring: %w", err)
	}
	return nil
}
