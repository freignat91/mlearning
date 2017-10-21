package mlapi

import (
	"github.com/freignat91/mlearning/mlserver/server"
	"golang.org/x/net/context"
)

//ServerLogs .
func (api *MlAPI) ServerLogs() ([]string, error) {
	client, _ := api.getClient()
	lines, err := client.client.ServerLogs(context.Background(),
		&mlserver.ServerLogsRequest{},
	)
	if err != nil {
		return nil, err
	}
	return lines.Lines, err
}

//SelectNetwork .
func (api *MlAPI) ServerSelectAnt(nestID int, antID int, mode string) error {
	client, _ := api.getClient()
	_, err := client.client.ServerSelectAnt(context.Background(),
		&mlserver.ServerSelectAntRequest{
			NestId: int32(nestID),
			AntId:  int32(antID),
			Mode:   mode,
		},
	)
	return err
}
