package exec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Exec(prog string, args ...string) (stdOut, stdErr bytes.Buffer, err error) {
	path, err := exec.LookPath(prog)
	if err != nil {
		err = fmt.Errorf("could not find '%s' executable in PATH. error: %w", prog, err)
		return
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err = cmd.Run()
	return
}

func DirExists(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("dir doesn't exist: %s, %w", path, err)
	}
	return nil
}
