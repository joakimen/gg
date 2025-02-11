package main

import (
	"context"

	"github.com/spf13/cobra"
)

func main() {
	cli := newRootCmd()
	ctx := cli.ExecuteContext(context.Background())
	cobra.CheckErr(ctx)
}
