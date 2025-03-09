// keyring package wraps the go-keyring dependency, and
// provides general methods for managing keyring items.

package keyring

import (
	"github.com/joakimen/gg"
	"github.com/zalando/go-keyring"
)

// Ensure Keyring implements KeyringProvider.
var _ gg.KeyringProvider = (*Keyring)(nil)

const keyringService = "gg"

type Keyring struct {
	service string
	user    string
}

func New(user string) *Keyring {
	return &Keyring{
		service: keyringService,
		user:    user,
	}
}

func (m *Keyring) Get() (string, error) {
	return keyring.Get(m.service, m.user)
}

func (m *Keyring) Set(val string) error {
	return keyring.Set(m.service, m.user, val)
}

func (m *Keyring) Delete() error {
	return keyring.Delete(m.service, m.user)
}

// KeyringProvider is the interface that wraps the basic keyring methods.
