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

// @title DID Resolver for did:cheqd method
// @version 1.0
// @description Universal Resolver driver for did:cheqd method

// @contact.name Cheqd Foundation Limited
// @contact.url https://cheqd.io

// @license.name Apache 2.0
// @license.url https://github.com/cheqd/did-resolver/blob/main/LICENSE

// @host resolver.cheqd.net
// @BasePath /1.0/identifiers
// @schemes https http

func main() {
	rootCmd := cmd.GetRootCmd()
	rootCmd.AddCommand(getVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
