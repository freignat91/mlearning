package main

import (
	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

type displayOptions struct {
	coef bool
}

var (
	displayOpts = displayOptions{}
)

//DisplayCmd .
var DisplayCmd = &cobra.Command{
	Use:   "display",
	Short: "display network",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.display(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	DisplayCmd.Flags().BoolVar(&displayOpts.coef, "coef", false, "display link coefs")
	NetworkCmd.AddCommand(DisplayCmd)
}

func (m *mlCLI) display(cmd *cobra.Command, args []string) error {
	api := mlapi.New(m.server)
	lines, err := api.Display(displayOpts.coef)
	if err != nil {
		return err
	}
	displayList(lines)
	return nil
}
