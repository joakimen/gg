package cli

import (
	"fmt"

	"github.com/joakimen/gg/internal/build"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(_ *cobra.Command, _ []string) {
			runVersion()
		},
	}
}

func runVersion() {
	fmt.Println(build.Version)
}
