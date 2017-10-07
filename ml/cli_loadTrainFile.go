package main

import (
	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

// LoadTrainFileCmd .
var LoadTrainFileCmd = &cobra.Command{
	Use:   "loadTrainFile",
	Short: "load json train file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.loadTrainFile(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NetworkCmd.AddCommand(LoadTrainFileCmd)
}

func (m *mlCLI) loadTrainFile(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		m.Fatal("usage: filePath\n")
	}
	api := mlapi.New(m.server)
	lines, err := api.LoadTrainFile(args[0])
	if err != nil {
		return err
	}
	displayList(lines)
	return nil
}
