package gg

type KeyringItemProvider interface {
	Get(string) (string, error)
	Set(string, string) error
	Delete(string) error
}
