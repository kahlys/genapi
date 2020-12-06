package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [project_name]",
	Short: "Create a new empty project",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// TODO: project only with letters - and _
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.ToLower(args[0])

		err := os.Mkdir(projectName, os.ModePerm)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}

		f, err := os.Create(filepath.Join(projectName, "config.yaml"))
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}
		_, err = f.WriteString(strings.TrimSpace(strings.Replace(configText, "PROJECT_NAME", projectName, -1)))
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}
	},
}

const configText = `ServiceName: PROJECT_NAME
Endpoints:
  - Name: "GetElem"
    Method: "GET"
    URL: "/api/elem/{id}"
  - Name: "SetElem"
    Method: "POST"
    URL: "/api/elem/{id}"
`
