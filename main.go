package main

import (
	"github.com/cheqd/did-resolver/cmd"
)

func main() {
	rootCmd := cmd.GetRootCmd()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
