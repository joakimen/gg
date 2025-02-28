package main

import (
	"fmt"
	"os"

	"github.com/joakimen/gg/cmd/gg/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
