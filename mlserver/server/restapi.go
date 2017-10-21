package mlserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

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

// FoodCoord .
type FoodCoord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//AntSelected .
type AntSelected struct {
	Nest int `json:"nest"`
	Ant  int `json:"ant"`
}

func (s *Server) getData(w http.ResponseWriter, r *http.Request) {
	data := s.nests.GetGraphicData()
	//verifJson(data)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) nestsStart(w http.ResponseWriter, r *http.Request) {
	s.nests.Start()
	saveInfo(r, "nestsStart")
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) nestsStop(w http.ResponseWriter, r *http.Request) {
	s.nests.Stop()
	saveInfo(r, "nestsStop")
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) isStarted(w http.ResponseWriter, r *http.Request) {
	ret := RetBool{Ret: s.nests.IsStarted()}
	json.NewEncoder(w).Encode(&ret)
}

func (s *Server) nextTime(w http.ResponseWriter, r *http.Request) {
	s.nests.NextTime()
	saveInfo(r, "nextTime")
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) exportAntSample(w http.ResponseWriter, r *http.Request) {
	nn, err := s.nests.ExportSelectedAntSample()
	if err != nil {
		fmt.Printf("Error exporting ant sample: %v\n", err)
		w.WriteHeader(400)
	}
	ret := RetInt{Ret: nn}
	saveInfo(r, "exportAntSample")
	json.NewEncoder(w).Encode(&ret)
}

func (s *Server) setSleep(w http.ResponseWriter, r *http.Request) {
	val, _ := strconv.Atoi(mux.Vars(r)["value"])
	s.nests.SetSleep(val)
	saveInfo(r, fmt.Sprintf("setSleep: %d", val))
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) setSelected(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var a AntSelected
	decoder.Decode(&a)
	s.nests.SetSelected(a.Nest, a.Ant, "")
	saveInfo(r, "setSelected")
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) getGlobalInfo(w http.ResponseWriter, r *http.Request) {
	saveInfo(r, "getGlobalInfo")
	info := s.nests.GetGlobalInfo()
	//verifJson(info)
	json.NewEncoder(w).Encode(info)
}

func (s *Server) getInfo(w http.ResponseWriter, r *http.Request) {
	info := s.nests.GetInfo()
	//verifJson(info)
	json.NewEncoder(w).Encode(info)
}

func (s *Server) restart(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var a RetInt
	decoder.Decode(&a)
	s.initNests(a.Ret)
	saveInfo(r, fmt.Sprintf("restart: %d", a.Ret))
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) addFoods(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t FoodCoord
	decoder.Decode(&t)
	s.nests.AddFoodGroup(t.X, t.Y)
	saveInfo(r, "addFood")
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) foodRenew(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t RetBool
	decoder.Decode(&t)
	s.nests.FoodRenew(t.Ret)
	saveInfo(r, "foodRenew")
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) setPanicMode(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t RetBool
	decoder.Decode(&t)
	s.nests.SetPanicMode(t.Ret)
	saveInfo(r, "panicMode")
	json.NewEncoder(w).Encode("{}")
}

func (s *Server) clearFoodGroup(w http.ResponseWriter, r *http.Request) {
	s.nests.ClearFoodGroup()
	saveInfo(r, "clearFoodGroup")
	json.NewEncoder(w).Encode("{}")
}

func saveInfo(r *http.Request, line string) {
	filename := "../addr.txt"
	addr := fmt.Sprintf("remote addr: %v\n", r.RemoteAddr)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("error opening in addr file: %v\n", err)
	}
	defer f.Close()
	if _, err = f.WriteString(addr + ": " + line); err != nil {
		fmt.Printf("error writing in addr file: %v\n", err)
	}
}

func verifJSON(v interface{}) {
	res, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("%s\n", res)
}
