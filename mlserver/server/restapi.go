package mlserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/freignat91/mlearning/nests"
	"github.com/gorilla/mux"
)

//RetBool .
type RetBool struct {
	Ret bool `json:"ret"`
}

//RetInt .
type RetInt struct {
	Ret int `json:"ret"`
}

func (s *Server) getData(w http.ResponseWriter, r *http.Request) {
	data := s.nests.GetData()
	json.NewEncoder(w).Encode(data)
}

func (s *Server) nestsStart(w http.ResponseWriter, r *http.Request) {
	s.nests.Start()
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) nestsStop(w http.ResponseWriter, r *http.Request) {
	s.nests.Stop()
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) isStarted(w http.ResponseWriter, r *http.Request) {
	ret := RetBool{Ret: s.nests.IsStarted()}
	json.NewEncoder(w).Encode(&ret)
}

func (s *Server) nextTime(w http.ResponseWriter, r *http.Request) {
	s.nests.NextTime()
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) exportAntSample(w http.ResponseWriter, r *http.Request) {
	nn, err := s.nests.ExportSelectedAntSample()
	if err != nil {
		fmt.Printf("Error exporting ant sample: %v\n", err)
		w.WriteHeader(400)
	}
	ret := RetInt{Ret: nn}
	json.NewEncoder(w).Encode(&ret)
}

func (s *Server) setSleep(w http.ResponseWriter, r *http.Request) {
	val, _ := strconv.Atoi(mux.Vars(r)["value"])
	s.nests.SetSleep(val)
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) setSelected(w http.ResponseWriter, r *http.Request) {
	val, _ := strconv.Atoi(mux.Vars(r)["selected"])
	s.nests.SetSelected(val)
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) getGlobalInfo(w http.ResponseWriter, r *http.Request) {
	info := s.nests.GetGlobalInfo()
	json.NewEncoder(w).Encode(info)
}

func (s *Server) getInfo(w http.ResponseWriter, r *http.Request) {
	info := s.nests.GetInfo()
	json.NewEncoder(w).Encode(info)
}

func (s *Server) restart(w http.ResponseWriter, r *http.Request) {
	nests, _ := nests.NewNests(0, 0, 500, 500, nestsDef)
	s.nests = nests
	json.NewEncoder(w).Encode("{}")
}