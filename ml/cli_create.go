package main

import (
	"fmt"
	"strconv"

	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

// CreateNetworkCmd .
var CreateNetworkCmd = &cobra.Command{
	Use:   "create",
	Short: "create network nbNeuroneLayer1, nbNeuroneLayer2,...",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.createNetwork(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NetworkCmd.AddCommand(CreateNetworkCmd)
}

func (m *mlCLI) createNetwork(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		m.Fatal("It needs at least in and out layers\n")
	}
	layers := make([]int32, len(args))
	for ii, snb := range args {
		nb, _ := strconv.Atoi(snb)
		layers[ii] = int32(nb)
	}
	api := mlapi.New(m.server)
	err := api.CreateNetwork(layers)
	if err != nil {
		return err
	}
	fmt.Println("network created")
	return nil
}
