package main

import (
	"github.com/cheqd/did-resolver/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	rootCmd := cmd.GetRootCmd()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
