package tty

import (
	"fmt"
	"syscall"

	"github.com/joakimen/gg"
	"golang.org/x/term"
)

var _ gg.TTYProvider = (*Provider)(nil)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (tty *Provider) Read(prompt string) (string, error) {
	fmt.Print(prompt)
	bytepw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()
	return string(bytepw), nil
}
