// keychain package wraps the go-keychain dependency
//
// This package should be usage agnostic, and provide simple Get/Set/Clear methods.

package keychain

import (
	"github.com/zalando/go-keyring"
)

type KeychainManager struct {
	service string
}

func (kcm *KeychainManager) Get(key string) (string, error) {
	return keyring.Get(kcm.service, key)
}

func (kcm *KeychainManager) Set(key string, val string) error {
	return keyring.Set(kcm.service, key, val)
}

func (kcm *KeychainManager) Delete(key string) error {
	return keyring.Delete(kcm.service, key)
}
