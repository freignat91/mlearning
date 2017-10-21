package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/freignat91/mlearning/api"
	"github.com/spf13/cobra"
)

type serverSelectAntOptions struct {
	mode string
}

var (
	serverSelectAntOpts = serverSelectAntOptions{}
)

// ServerSelectAntCmd .
var ServerSelectAntCmd = &cobra.Command{
	Use:   "selectAnt",
	Short: "select an Ant and set its network as current of a given running ant: nestId, antId or workerBest or soldierBest or selected",
	Run: func(cmd *cobra.Command, args []string) {
		if err := mlCli.serverSelectAnt(cmd, args); err != nil {
			mlCli.Fatal("Error: %v\n", err)
		}
	},
}

func init() {
	ServerSelectAntCmd.Flags().StringVar(&serverSelectAntOpts.mode, "what", "", "bestWorker or bestSoldier or selected")
	ServerCmd.AddCommand(ServerSelectAntCmd)
}

func (m *mlCLI) serverSelectAnt(cmd *cobra.Command, args []string) error {
	nestID := 0
	antID := 0
	if strings.ToLower(serverSelectAntOpts.mode) != "selected" {
		if len(args) >= 1 {
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
	if err := api.ServerSelectAnt(nestID, antID, serverSelectAntOpts.mode); err != nil {
		return err
	}
	if nestID == 0 && antID == 0 {
		fmt.Println("Selected ant network selected")
	} else {
		fmt.Printf("Network %d-%d selected\n", nestID, antID)
	}
	return nil
}
