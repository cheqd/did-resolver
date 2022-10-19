package main

import (
	"github.com/cheqd/did-resolver/cmd"
	_ "github.com/cheqd/did-resolver/docs"
	"github.com/spf13/cobra"
)

var version = "dev"

func getVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of the binary",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printVersion(cmd, args)
		},
	}
}

func printVersion(cmd *cobra.Command, args []string) error {
	println(version)

	return nil
}

// @title Cheqd DID Resolver API
// @version 1.0
// @description Cheqd DID Resolver API for DID resolution and dereferencing.

// @contact.name Cheqd
// @contact.url https://cheqd.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	rootCmd := cmd.GetRootCmd()
	rootCmd.AddCommand(getVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
