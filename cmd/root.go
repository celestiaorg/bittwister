package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func Execute() {
	if len(os.Args) > 0 {
		rootCmd.Use = os.Args[0]
	}

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {

		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
