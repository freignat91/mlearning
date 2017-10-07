package main

import (
	"fmt"
	"strconv"

	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

// PropagateCmd .
var PropagateCmd = &cobra.Command{
	Use:   "propagate",
	Short: "push value to input layer value1, value2, ...",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.propagate(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NetworkCmd.AddCommand(PropagateCmd)
}

func (m *mlCLI) propagate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		m.Fatal("at lest one argument is mandatory\n")
	}
	values := make([]float64, 0)
	for ii := 0; ii < len(args); ii++ {
		value, _ := strconv.ParseFloat(args[ii], 64)
		values = append(values, value)
	}
	api := mlapi.New(m.server)
	outs, err := api.Propagate(values)
	if err != nil {
		return err
	}
	fmt.Printf("Outs: %v\n", outs)
	return nil
}
