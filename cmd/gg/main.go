package main

import (
	"context"

	"github.com/joakimen/gg/pkg/cmd/root"
	"github.com/spf13/cobra"
)

func main() {
	cli := root.NewCmd()
	ctx := cli.ExecuteContext(context.Background())
	cobra.CheckErr(ctx)
}
