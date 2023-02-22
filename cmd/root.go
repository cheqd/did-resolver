package cmd

import (
	"github.com/spf13/cobra"
)

func GetRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "did-resolver",
		Short: "DID resolver for the cheqd method",
	}

	rootCmd.AddCommand(getServeCmd(), getPrintConfigCmd(), getVersionCmd())

	return rootCmd
}
