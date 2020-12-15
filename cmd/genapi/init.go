package main

import (
	"fmt"
	"os"

	"github.com/kahlys/genapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init an empty project",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInitProject(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}
		if err := os.Remove("config.yaml"); err != nil {
			fmt.Printf("WARNING: %v\n", err)
			return
		}
	},
}

var restapi genapi.RestAPI

func runInitProject() error {
	if err := parseConfig(); err != nil {
		return err
	}
	if err := restapi.Generate(); err != nil {
		return err
	}
	return nil
}

func parseConfig() error {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("unable to read configuration file: %v", err)
	}
	// Configuration for the micro-service
	if err := viper.Unmarshal(&restapi); err != nil {
		return fmt.Errorf("unable to decode into config struct: %v", err)
	}
	return nil
}
