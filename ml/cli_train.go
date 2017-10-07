package main

import (
	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

type trainOptions struct {
	number        int
	createNetwork bool
	all           bool
	hide          bool
	analyse       bool
}

var (
	trainOpts = trainOptions{}
)

// TrainCmd .
var TrainCmd = &cobra.Command{
	Use:   "train",
	Short: "train network using logical train file name",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.train(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	TrainCmd.Flags().IntVarP(&trainOpts.number, "number", "n", 1, "number of train")
	TrainCmd.Flags().BoolVarP(&trainOpts.createNetwork, "createNetwork", "c", false, "create the network before train")
	TrainCmd.Flags().BoolVarP(&trainOpts.all, "all", "a", false, "train the whole data -n (--number) time")
	TrainCmd.Flags().BoolVar(&trainOpts.hide, "hide", false, "hide sample details")
	TrainCmd.Flags().BoolVar(&trainOpts.analyse, "analyse", false, "display an analyse of the sample")
	NetworkCmd.AddCommand(TrainCmd)
}

func (m *mlCLI) train(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		m.Fatal("need logical train file name as first argument\n")
	}
	api := mlapi.New(m.server)
	lines, err := api.Train(args[0], trainOpts.number, trainOpts.all, trainOpts.hide, trainOpts.createNetwork, trainOpts.analyse)
	if err != nil {
		return err
	}
	displayList(lines)
	return nil
}
