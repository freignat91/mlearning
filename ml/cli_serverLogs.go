package main

import (
	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

// ServerLogsCmd .
var ServerLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "server logs toogles",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.serverLogs(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	ServerCmd.AddCommand(ServerLogsCmd)
}

func (m *mlCLI) serverLogs(cmd *cobra.Command, args []string) error {
	api := mlapi.New(m.server)
	lines, err := api.ServerLogs()
	if err != nil {
		return err
	}
	displayList(lines)
	return nil
}
