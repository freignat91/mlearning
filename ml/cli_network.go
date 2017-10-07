package main

import (
	"github.com/spf13/cobra"
)

// NetworkCmd .
var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "network test operations",
	//Aliases: []string{"pf"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(NetworkCmd)
}
