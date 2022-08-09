package main

import (
	"github.com/cheqd/did-resolver/cmd"
	"github.com/spf13/cobra"
)

var (
    Version = "dev"
    Commit  = "none"
    Date    = "unknown"
    BuiltBy = "unknown"
)

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
	println(Version)

	return nil
}

func main() {
	rootCmd := cmd.GetRootCmd()
	rootCmd.AddCommand(getVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
