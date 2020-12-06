package main

import (
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fcosctl",
		Short: "Fast Fedora CoreOS configuration development tool",
	}
}

func main() {
	cmd := newRootCmd()

	cmd.AddCommand(newRunCmd())
	cmd.AddCommand(newImageCmd())

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
