package main

import (
	"fmt"
	"strconv"

	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

// TestNetworkCmd .
var TestNetworkCmd = &cobra.Command{
	Use:   "test",
	Short: "test the network of a given running ant: nestId, antId",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.testNetwork(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NetworkCmd.AddCommand(TestNetworkCmd)
}

func (m *mlCLI) testNetwork(cmd *cobra.Command, args []string) error {
	nestID := 0
	antID := 0
	if len(args) >= 1 {
		if args[0] == "best" {
			nestID = -1
		} else if args[0] == "worse" {
			nestID = -2
		} else {
			nid, err1 := strconv.Atoi(args[0])
			if err1 != nil {
				return fmt.Errorf("nestId is not a number: %s", args[0])
			}
			nestID = nid
		}
		if len(args) == 2 {
			aid, err2 := strconv.Atoi(args[1])
			if err2 != nil {
				return fmt.Errorf("antId is not a number: %s", args[1])
			}
			antID = aid
		}
	}
	api := mlapi.New(m.server)
	lines, err := api.TestNetwork(nestID, antID)
	if err != nil {
		return err
	}
	displayList(lines)
	return nil
}
