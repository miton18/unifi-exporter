package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"info"},
		Short:   "Show build informations of Enifi exporter",
		Long:    "Display build commit, date, and version of Unifi exporter",
		Run:     versionCmdFn,
	}

	format string
)

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")
}

func versionCmdFn(cmd *cobra.Command, args []string) {
	switch format {
	case "json":
		txt, _ := json.MarshalIndent(map[string]interface{}{
			"name":    "Unifi exporter",
			"version": version,
			"build": map[string]string{
				"commit": commit,
				"date":   date,
			},
		}, "", " ")
		fmt.Println(string(txt))

	default:
		fmt.Println("Unifi exporter")
		fmt.Println("==============")
		fmt.Println("Version: " + version)
		fmt.Println("Git commit: " + commit)
		fmt.Println("Build date: " + date)
	}
}
