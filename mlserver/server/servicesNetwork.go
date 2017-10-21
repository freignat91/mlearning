package mlserver

import (
	"fmt"

	"github.com/freignat91/mlearning/network"

	"golang.org/x/net/context"
)

//CreateNetwork create new network
func (s *Server) CreateNetwork(ctx context.Context, req *CreateNetworkRequest) (*EmptyReply, error) {
	layers := make([]int, len(req.Layers))
	for ii, val := range req.Layers {
		layers[ii] = int(val)
	}
	net, err := network.NewNetwork(layers)
	if err != nil {
		return nil, err
	}
	s.network = net
	return &EmptyReply{}, nil
}

//Propagate .
func (s *Server) Propagate(ctx context.Context, req *PropagateRequest) (*PropagateReply, error) {
	if s.network == nil {
		return nil, fmt.Errorf("the network is not created")
	}
	fmt.Printf("Execute: Propagate\n")
	outs := s.network.Propagate(req.InValues, true)
	return &PropagateReply{OutValues: outs}, nil
}

//BackPropagate .
func (s *Server) BackPropagate(ctx context.Context, req *BackPropagateRequest) (*EmptyReply, error) {
	if s.network == nil {
		return nil, fmt.Errorf("the network is not created")
	}
	fmt.Printf("Execute: BackPropagate\n")
	s.network.BackPropagate(req.OutValues)
	s.network.Display(true)
	return &EmptyReply{}, nil
}

//Display .
func (s *Server) Display(ctx context.Context, req *DisplayRequest) (*LinesReply, error) {
	if s.network == nil {
		return nil, fmt.Errorf("the network is not created")
	}
	lines := s.network.Display(req.Coef)
	return &LinesReply{Lines: lines}, nil
}

//LoadTrainFile .
func (s *Server) LoadTrainFile(ctx context.Context, req *LoadTrainFileRequest) (*LinesReply, error) {
	lines, err := s.network.LoadTrainFile(req.Path)
	return &LinesReply{Lines: lines}, err
}

//Train .
func (s *Server) Train(ctx context.Context, req *TrainRequest) (*LinesReply, error) {
	if req.CreateNetwork || !s.network.IsCreated() {
		n, err := network.NetNetworkFromDataSet(req.Name)
		if err != nil {
			return nil, err
		}
		s.network = n
	}
	lines, err := s.network.Train(req.Name, int(req.Number), req.All, req.Hide, req.CreateNetwork, req.Analyse)
	if err != nil {
		return nil, err
	}
	return &LinesReply{Lines: lines}, nil
}

//TrainSoluce .
func (s *Server) TrainSoluce(ctx context.Context, req *TrainSoluceRequest) (*LinesReply, error) {
	lines := s.nests.TrainSoluce(int(req.Number))
	return &LinesReply{Lines: lines}, nil
}

//Test .
func (s *Server) Test(ctx context.Context, req *TestRequest) (*LinesReply, error) {
	lines := s.nests.Test()
	return &LinesReply{Lines: lines}, nil
}
