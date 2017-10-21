package main

import (
	"strconv"

	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

type trainSoluceOptions struct {
}

var (
	trainSoluceOpts = trainSoluceOptions{}
)

// TrainSoluceCmd .
var TrainSoluceCmd = &cobra.Command{
	Use:   "trainSoluce",
	Short: "train network right computed samples",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.trainSoluce(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	NetworkCmd.AddCommand(TrainSoluceCmd)
}

func (m *mlCLI) trainSoluce(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		m.Fatal("need number of training as first argument\n")
	}
	nn, errc := strconv.Atoi(args[0])
	if errc != nil {
		m.Fatal("need number as first argument: %s\n", args[0])
	}
	api := mlapi.New(m.server)
	lines, err := api.TrainSoluce(nn)
	if err != nil {
		return err
	}
	displayList(lines)
	return nil
}
