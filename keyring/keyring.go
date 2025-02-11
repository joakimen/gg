// keyring package wraps the go-keyring dependency, and
// provides general methods for managing keyring items.

package keyring

import (
	"github.com/zalando/go-keyring"
)

type Manager struct {
	service string
}

func NewKeyringManager(service string) *Manager {
	return &Manager{
		service: service,
	}
}

func (m *Manager) Get(user string) (string, error) {
	return keyring.Get(m.service, user)
}

func (m *Manager) Set(user string, val string) error {
	return keyring.Set(m.service, user, val)
}

func (m *Manager) Delete(user string) error {
	return keyring.Delete(m.service, user)
}
