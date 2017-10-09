package mlserver

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/freignat91/mlearning/nests"
	"github.com/freignat91/mlearning/network"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/urfave/negroni"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	config = mlConfig{}
	ctx    = context.Background()
)

const port = "3001"

//Server .
type Server struct {
	host    string
	conn    *grpc.ClientConn
	network *network.MLNetwork
	nests   *nests.Nests
}

var nestsDef = []int{100}

// Start gnode
func (s *Server) Start(version string) error {
	config.init(version)
	s.init()
	for {
		//
		time.Sleep(3000 * time.Second)
	}

}

func (s *Server) init() {
	s.network = &network.MLNetwork{}
	s.initNests()
	//fmt.Printf("init data: %v\n", s.nests.GetData())
	s.startGRPCServer()
	s.start()
}

func (s *Server) initNests() {
	nests, _ := nests.NewNests(0, 0, 500, 500, nestsDef, 50, 5)
	s.nests = nests
}

func (s *Server) start() {
	r := mux.NewRouter()
	n := negroni.Classic()
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(r)

	abspath, err := filepath.Abs("./public")
	if err != nil {
		fmt.Print(err)
	}
	s.handleAPIFunctions(r)
	fs := http.FileServer(http.Dir(abspath))
	r.PathPrefix("/").Handler(fs)
	log.Printf("Rest server starting on %s\n", port)
	if err := http.ListenAndServe(":"+port, n); err != nil {
		log.Fatal("Server error: ", err)
	}
}

func (s *Server) handleAPIFunctions(r *mux.Router) {
	r.HandleFunc("/api/v1/data", s.getData).Methods("GET")
	r.HandleFunc("/api/v1/start", s.nestsStart).Methods("GET")
	r.HandleFunc("/api/v1/stop", s.nestsStop).Methods("GET")
	r.HandleFunc("/api/v1/isStarted", s.isStarted).Methods("GET")
	r.HandleFunc("/api/v1/nextTime", s.nextTime).Methods("GET")
	r.HandleFunc("/api/v1/exportAntSample", s.exportAntSample).Methods("GET")
	r.HandleFunc("/api/v1/setSleep/{value}", s.setSleep).Methods("GET")
	r.HandleFunc("/api/v1/setSelected/{selected}", s.setSelected).Methods("GET")
	r.HandleFunc("/api/v1/globalInfo", s.getGlobalInfo).Methods("GET")
	r.HandleFunc("/api/v1/info", s.getInfo).Methods("GET")
	r.HandleFunc("/api/v1/restart", s.restart).Methods("GET")
	r.HandleFunc("/api/v1/addFoods", s.addFoods).Methods("POST")
	r.HandleFunc("/api/v1/foodRenew", s.foodRenew).Methods("POST")
	r.HandleFunc("/api/v1/clearFoodGroup", s.clearFoodGroup).Methods("POST")
}

//startGRPCServer .
func (s *Server) startGRPCServer() {
	sr := grpc.NewServer()
	RegisterMLearningServiceServer(sr, s)
	go func() {
		lis, err := net.Listen("tcp", ":"+config.grpcPort)
		if err != nil {
			fmt.Printf("mLearning is unable to listen on: %s\n%v", ":"+config.grpcPort, err)
		}
		fmt.Printf("mLearning is listening on port %s\n", ":"+config.grpcPort)
		if err := sr.Serve(lis); err != nil {
			fmt.Printf("Problem in mLearning server: %s\n", err)
		}
	}()
}
