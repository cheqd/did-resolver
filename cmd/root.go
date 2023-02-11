package cmd

import (
	"github.com/spf13/cobra"
)

func GetRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cheqd-did-resolver",
		Short: "DID resolver for the cheqd method",
	}

	rootCmd.AddCommand(getServeCmd(), getPrintConfigCmd())

	return rootCmd
}
