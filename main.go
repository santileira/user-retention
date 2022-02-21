package main

import (
	"github.com/santileira/user-retention/cmd"
	"os"
)

func main() {
	if err := cmd.Cmds().Execute(); err != nil {
		os.Exit(1)
	}
}
