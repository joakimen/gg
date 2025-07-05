package prompt

import (
	"fmt"
	"syscall"

	"github.com/joakimen/gg"
	"golang.org/x/term"
)

var _ gg.InputReader = ReadPassword

func ReadPassword() (string, error) {
	bytepw, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	return string(bytepw), nil
}
