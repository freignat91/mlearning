package main

import (
	"github.com/spf13/cobra"
)

// ServerCmd .
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "server operations",
	//Aliases: []string{"pf"},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	RootCmd.AddCommand(ServerCmd)
}
