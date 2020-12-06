package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "genapi",
		Short: "A generator for REST-API Applications",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(
		createCmd,
		initCmd,
	)

	rootCmd.Execute()
}
