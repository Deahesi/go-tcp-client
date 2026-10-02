package cmd

import "github.com/spf13/cobra"

func Execute() error {
	root := &cobra.Command{
		Use:           "client",
		Short:         "Client for sending files to the lab server",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.AddCommand(newSendCommand(), newGenerateCommand())
	return root.Execute()
}
