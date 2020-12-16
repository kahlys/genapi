package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "genapi",
		Short: "A generator for REST-API Applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(
		createCmd,
		initCmd,
	)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("FAILED: ", err)
		os.Exit(1)
	}
}
