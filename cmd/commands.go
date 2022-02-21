package cmd

import (
	"github.com/santileira/user-retention/cmd/script"
	"github.com/spf13/cobra"
)

func Cmds() *cobra.Command {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(script.NewRunnable().Cmd())
	return rootCmd
}
