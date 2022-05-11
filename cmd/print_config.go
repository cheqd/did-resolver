package cmd

import (
	"github.com/cheqd/cheqd-did-resolver/utils"
	"github.com/spf13/cobra"
)

func getPrintConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "print-config",
		Short: "Prints the active configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printConfig()
		},
	}
}

func printConfig() error {
	config := utils.MustLoadConfig()
	configYaml := config.MustMarshalYaml()

	println(configYaml)

	return nil
}
