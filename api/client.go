package mlapi

import (
	"time"

	"github.com/freignat91/mlearning/mlserver/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type mlClient struct {
	api      *MlAPI
	client   mlserver.MLearningServiceClient
	nodeHost string
	ctx      context.Context
	conn     *grpc.ClientConn
}

func (g *mlClient) init(api *MlAPI) error {
	g.api = api
	g.ctx = context.Background()
	if err := g.connectServer(); err != nil {
		return err
	}
	api.info("Client connected to server %s\n", g.nodeHost)
	return nil
}

func (g *mlClient) connectServer() error {
	cn, err := grpc.Dial("localhost:30107",
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*20))
	if err != nil {
		return err
	}
	g.conn = cn
	g.client = mlserver.NewMLearningServiceClient(g.conn)
	return nil
}

func (g *mlClient) close() {
	if g.conn != nil {
		g.conn.Close()
	}
}
