package main

import (
	"strconv"

	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

// BackPropagateCmd .
var BackPropagateCmd = &cobra.Command{
	Use:   "backPropagate",
	Short: "push value to out layer value1, value2, ...",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.backPropagate(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NetworkCmd.AddCommand(BackPropagateCmd)
}

func (m *mlCLI) backPropagate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		m.Fatal("at lest one argument is mandatory\n")
	}
	values := make([]float64, 0)
	for ii := 0; ii < len(args); ii++ {
		value, _ := strconv.ParseFloat(args[ii], 64)
		values = append(values, value)
	}
	api := mlapi.New(m.server)
	return api.BackPropagate(values)
}
