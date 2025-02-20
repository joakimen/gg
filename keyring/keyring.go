// keyring package wraps the go-keyring dependency, and
// provides general methods for managing keyring items.

package keyring

import (
	"github.com/zalando/go-keyring"
)

type Manager struct {
	service string
	user    string
}

func NewManager(service string, user string) *Manager {
	return &Manager{
		service: service,
		user:    user,
	}
}

func (m *Manager) Get() (string, error) {
	return keyring.Get(m.service, m.user)
}

func (m *Manager) Set(val string) error {
	return keyring.Set(m.service, m.user, val)
}

func (m *Manager) Delete() error {
	return keyring.Delete(m.service, m.user)
}
