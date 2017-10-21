package mlapi

import (
	"github.com/freignat91/mlearning/mlserver/server"
	"golang.org/x/net/context"
)

//CreateNetwork create a new network
func (api *MlAPI) CreateNetwork(layers []int32) error {
	client, _ := api.getClient()
	_, err := client.client.CreateNetwork(context.Background(),
		&mlserver.CreateNetworkRequest{
			Layers: layers,
		},
	)
	return err
}

//Propagate propagate from input to output layer
func (api *MlAPI) Propagate(values []float64) ([]float64, error) {
	client, _ := api.getClient()
	ret, err := client.client.Propagate(context.Background(),
		&mlserver.PropagateRequest{
			InValues: values,
		},
	)
	if err != nil {
		return nil, err
	}
	return ret.OutValues, nil
}

//BackPropagate propagate from input to output layer
func (api *MlAPI) BackPropagate(values []float64) error {
	client, _ := api.getClient()
	_, err := client.client.BackPropagate(context.Background(),
		&mlserver.BackPropagateRequest{
			OutValues: values,
		},
	)
	return err
}

//Display .
func (api *MlAPI) Display(coef bool) ([]string, error) {
	client, _ := api.getClient()
	lines, err := client.client.Display(context.Background(),
		&mlserver.DisplayRequest{
			Coef: coef,
		},
	)
	if err != nil {
		return nil, err
	}
	return lines.Lines, nil
}

//LoadTrainFile .
func (api *MlAPI) LoadTrainFile(path string) ([]string, error) {
	client, _ := api.getClient()
	lines, err := client.client.LoadTrainFile(context.Background(),
		&mlserver.LoadTrainFileRequest{
			Path: path,
		},
	)
	if err != nil {
		return nil, err
	}
	return lines.Lines, err
}

//Train .
func (api *MlAPI) Train(name string, nb int, all bool, hide bool, createNetwork bool, analyse bool) ([]string, error) {
	client, _ := api.getClient()
	lines, err := client.client.Train(context.Background(),
		&mlserver.TrainRequest{
			Name:          name,
			Number:        int32(nb),
			All:           all,
			Hide:          hide,
			Analyse:       analyse,
			CreateNetwork: createNetwork,
		},
	)
	if err != nil {
		return nil, err
	}
	return lines.Lines, err
}

//TrainSoluce .
func (api *MlAPI) TrainSoluce(nb int) ([]string, error) {
	client, _ := api.getClient()
	lines, err := client.client.TrainSoluce(context.Background(),
		&mlserver.TrainSoluceRequest{
			Number: int32(nb),
		},
	)
	if err != nil {
		return nil, err
	}
	return lines.Lines, err
}

//TestNetwork .
func (api *MlAPI) TestNetwork(nestID int, antID int) ([]string, error) {
	client, _ := api.getClient()
	lines, err := client.client.Test(context.Background(),
		&mlserver.TestRequest{
			NestId: int32(nestID),
			AntId:  int32(antID),
		},
	)
	if err != nil {
		return nil, err
	}
	return lines.Lines, nil
}
