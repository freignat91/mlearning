package mlserver

import "golang.org/x/net/context"

//ServerLogs .
func (s *Server) ServerLogs(ctx context.Context, req *ServerLogsRequest) (*LinesReply, error) {
	lines := s.nests.LogsToggle()
	return &LinesReply{Lines: lines}, nil
}

//ServerSelectAnt .
func (s *Server) ServerSelectAnt(ctx context.Context, req *ServerSelectAntRequest) (*EmptyReply, error) {
	s.nests.SetSelected(int(req.NestId), int(req.AntId), req.Mode)
	network, err := s.nests.GetSelectedNetwork()
	if err != nil {
		return nil, err
	}
	s.network = network
	return &EmptyReply{}, nil
}
