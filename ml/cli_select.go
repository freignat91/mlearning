package main

import (
	"fmt"
	"strconv"

	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

// SelectNetworkCmd .
var SelectNetworkCmd = &cobra.Command{
	Use:   "select",
	Short: "select the network of a given running ant: nestId, antId",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.selectNetwork(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NetworkCmd.AddCommand(SelectNetworkCmd)
}

func (m *mlCLI) selectNetwork(cmd *cobra.Command, args []string) error {
	nestID := 0
	antID := 0
	if len(args) == 2 {
		nid, err1 := strconv.Atoi(args[0])
		if err1 != nil {
			return fmt.Errorf("nestId is not a number: %s", args[0])
		}
		nestID = nid
		aid, err2 := strconv.Atoi(args[1])
		if err2 != nil {
			return fmt.Errorf("antId is not a number: %s", args[1])
		}
		antID = aid
	}
	api := mlapi.New(m.server)
	if err := api.SelectNetwork(nestID, antID); err != nil {
		return err
	}
	if nestID == 0 && antID == 0 {
		fmt.Println("Selected ant network selected")
	} else {
		fmt.Printf("Network %d-%d selected\n", nestID, antID)
	}
	return nil
}
