package tty

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

type PasswordProvider struct{}

func NewPasswordProvider() *PasswordProvider {
	return &PasswordProvider{}
}

func (pp *PasswordProvider) ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytepw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println()
	return string(bytepw), nil
}
