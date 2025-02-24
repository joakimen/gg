package version

import (
	"fmt"

	"github.com/joakimen/gg/internal/build"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println(build.Version)
			return nil
		},
	}

	return cmd
}
