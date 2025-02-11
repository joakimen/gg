package github

import "github.com/joakimen/gg/keyring"

type CredentialsManager struct {
	keyringManager keyring.Manager
	keyringUser    string
}

func NewCredentialsManager() CredentialsManager {
	const (
		keyringService = "gg"
		keyringUser    = "github"
	)
	return CredentialsManager{
		keyringManager: *keyring.NewKeyringManager(keyringService),
		keyringUser:    keyringUser,
	}
}

func (cm *CredentialsManager) GetToken() (string, error) {
	return cm.keyringManager.Get(cm.keyringUser)
}

func (cm *CredentialsManager) SetToken(token string) error {
	return cm.keyringManager.Set(cm.keyringUser, token)
}

func (cm *CredentialsManager) ClearToken() error {
	return cm.keyringManager.Delete(cm.keyringUser)
}
