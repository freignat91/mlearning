package nests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/freignat91/mlearning/network"
)

//Nests .
type Nests struct {
	waiter               int
	nbs                  []int
	nests                []*Nest
	totalNumber          int64
	attractors           *Attractors
	xmin                 float64
	xmax                 float64
	ymin                 float64
	ymax                 float64
	stopped              bool
	timeRef              int64
	lastTimeRef          int64
	speed                int64
	random               bool
	selectedNest         int
	selected             int
	averageRate          float64
	bestNest             *Nest
	worseNest            *Nest
	bestAnt              *Ant
	worseAnt             *Ant
	ready                bool
	happiness            float64
	dataSet              *network.MlDataSet
	foods                []*Food
	foodGroups           []*FoodGroup
	foodRenew            bool
	pheromones           []*Pheromone
	pheromoneLevel       float64
	pheromoneAntDelay    int
	pheromoneGroup       float64
	pheromoneFadeDelay   int
	pheromoneFadeCounter int
	//
	log           bool
	period        int64
	statTrain     *Stats
	statDecision  *Stats
	statReinforce *Stats
	statFade      *Stats
	statNetwork   *Stats
	statContact   *Stats
	statFood      *Stats
}

// GraphicData .
type GraphicData struct {
	Ants       []*AntData   `json:"ants"`
	Foods      []*Food      `json:"foods"`
	Pheromones []*Pheromone `json:"pheromones"`
}

//AntData .
type AntData struct {
	ID        int     `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Direction int     `json:"direction"`
	Contact   bool    `json:"contact"`
	//Entries   []float64 `json:"entries"`
}

//GlobalInfo .
type GlobalInfo struct {
	Nests        []int   `json:"nests"`
	Xmin         float64 `json:"xmin"`
	Xmax         float64 `json:"xmax"`
	Ymin         float64 `json:"ymin"`
	Ymax         float64 `json:"ymax"`
	Ndir         int     `json:"ndir"`
	Waiter       int     `json:"waiter"`
	SelectedNest int     `json:"selectedNest"`
	SelectedAnt  int     `json:"selectedAnt"`
}

//Infos .
type Infos struct {
	Timer                int64   `json:"timer"`
	Speed                int64   `json:"speed"`
	Selected             *Info   `json:"selected"`
	Global               *Info   `json:"global"`
	BestID               int     `json:"bestId"`
	BestNetworkStruct    string  `json:"bestNetworkStruct"`
	BestNetworkGRate     float64 `json:"bestNetworkGRate"`
	BestNetworkDirCount  int     `json:"bestNetworkDirCount"`
	WorseID              int     `json:"worseId"`
	WorseNetworkStruct   string  `json:"worseNetworkStruct"`
	WorseNetworkGRate    float64 `json:"worseNetworkGRate"`
	WorseNetworkDirCount int     `json:"worseNetworkDirCount"`
	FromBeginningFoods   int64   `json:"fromBeginningFoods"`
	PeriodFoods          int64   `json:"periodFoods"`
}

//Info .
type Info struct {
	Happiness                   float64 `json:"happiness"`
	Train                       int64   `json:"train"`
	FromBeginningTrain          int64   `json:"fromBeginningTrain"`
	Reinforce                   int64   `json:"reinforce"`
	Fade                        int64   `json:"fade"`
	Decision                    int64   `json:"decision"`
	Contact                     int64   `json:"contact"`
	PeriodNetworkUpdated        int64   `json:"periodNetworkUpdated"`
	PeriodGRate                 float64 `json:"periodGRate"`
	FromBeginningNetworkUpdated int64   `json:"fromBeginningNetworkUpdated"`
	FromBeginningGRate          float64 `json:"fromBeginningGRate"`
}

//NewNests .
func NewNests(xmin float64, ymin float64, xmax float64, ymax float64, nbs []int, foodNb int, foodGroupNb int) (*Nests, error) {
	nests := &Nests{
		xmin:               xmin,
		ymin:               ymin,
		xmax:               xmax,
		ymax:               ymax,
		waiter:             1,
		nbs:                nbs,
		nests:              make([]*Nest, len(nbs), len(nbs)),
		stopped:            true,
		random:             false, //to compate with random ant behavior
		selected:           0,
		selectedNest:       1,
		period:             10000,
		log:                false,
		foods:              make([]*Food, 0, 0),
		foodGroups:         make([]*FoodGroup, 0, 0),
		foodRenew:          true,
		pheromones:         make([]*Pheromone, 0, 0),
		pheromoneLevel:     1000,
		pheromoneAntDelay:  20,
		pheromoneFadeDelay: 5,
		pheromoneGroup:     100,
		statTrain:          newStats(nil, nil),
		statDecision:       newStats(nil, nil),
		statReinforce:      newStats(nil, nil),
		statFade:           newStats(nil, nil),
		statNetwork:        newStats(nil, nil),
		statContact:        newStats(nil, nil),
		statFood:           newStats(nil, nil),
	}
	for _, nb := range nbs {
		nests.totalNumber += int64(nb)
	}
	for ii, nb := range nbs {
		nest, err := newNest(nests, ii+1, nb)
		if err != nil {
			return nil, err
		}
		nests.nests[ii] = nest
	}
	nests.bestNest = nests.nests[0]
	nests.worseNest = nests.nests[0]
	nests.bestAnt = nests.bestNest.bestAnt
	nests.worseAnt = nests.worseNest.worseAnt
	nests.attractors = newAttractors()
	nests.ready = true
	return nests, nil
}

//IsReady .
func (ns *Nests) IsReady() bool {
	return ns.ready
}

//GetGraphicData .
func (ns *Nests) GetGraphicData() *GraphicData {
	ants := make([]*AntData, 0, 0)
	for _, nest := range ns.nests {
		nest.addData(&ants)
	}
	return &GraphicData{
		Ants:       ants,
		Foods:      ns.foods,
		Pheromones: ns.pheromones,
	}
}

//Start .
func (ns *Nests) Start() {
	if !ns.stopped {
		return
	}
	ns.stopped = false
	go func() {
		for {
			ns.NextTime()
			if ns.waiter > 0 {
				time.Sleep(time.Duration(ns.waiter) * time.Millisecond)
			}
			if ns.stopped {
				return
			}
		}
	}()
}

//Stop .
func (ns *Nests) Stop() {
	if ns.selected > 0 {
		aa := ns.nests[ns.selectedNest-1].ants[ns.selected-1]
		aa.commitPeriodStats(ns)
	}
	ns.stopped = true
}

//IsStarted .
func (ns *Nests) IsStarted() bool {
	return !ns.stopped
}

//NextTime .
func (ns *Nests) NextTime() {
	//ns.printf("\033[2J\033[0;0H\n")
	ns.timeRef++
	h := 0.0
	for _, nest := range ns.nests {
		nest.nextTime(ns)
		h += nest.happiness
		ns.fadePheromones()
		if nest.bestAnt.gRate > ns.bestNest.bestAnt.gRate {
			ns.bestNest = nest
		}
		if nest.worseAnt.gRate < ns.worseNest.worseAnt.gRate {
			ns.worseNest = nest
		}
	}
	ns.bestAnt = ns.bestNest.bestAnt
	ns.worseAnt = ns.worseNest.worseAnt
	ns.averageRate = float64(ns.statReinforce.cumul) * 100.0 / float64(ns.statDecision.cumul)
	ns.happiness = h / float64(len(ns.nests))
}

func (ns *Nests) printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (ns *Nests) getSelected() *Ant {
	if ns.selected <= 0 {
		return nil
	}
	return ns.nests[ns.selectedNest-1].ants[ns.selected-1]
}

func (ns *Nests) commitPeriodStats() {
	if ns.stopped {
		return
	}
	ns.speed = ns.timeRef - ns.lastTimeRef
	ns.lastTimeRef = ns.timeRef
	for _, nest := range ns.nests {
		nest.statTrain.push()
		nest.statDecision.push()
		nest.statReinforce.push()
		nest.statFade.push()
		nest.statNetwork.push()
		nest.statContact.push()
		nest.statFood.push()
	}
	ns.statTrain.push()
	ns.statDecision.push()
	ns.statReinforce.push()
	ns.statFade.push()
	ns.statNetwork.push()
	ns.statContact.push()
	ns.statFood.push()
}

//SetSelected .
func (ns *Nests) SetSelected(selected int) {
	ns.printf("selected: %d\n", selected)
	ns.selected = selected
	ant := ns.nests[ns.selectedNest-1].ants[ns.selected-1]
	ns.dataSet = &network.MlDataSet{
		Name:   "ant",
		Layers: ant.network.Getdef(),
		Data:   make([]network.MlDataSample, 0, 0),
	}
	if ns.stopped {
		ant.commitPeriodStats(ns)
	}
}

func (ns *Nests) addSample(ins []float64, outs []float64) {
	if len(ns.dataSet.Data) >= 10000 {
		return
	}
	mouts := make([]float64, outNb, outNb)
	for ii, out := range outs {
		mouts[ii] = out
	}
	data := network.MlDataSample{
		In:  ins,
		Out: mouts,
	}
	//a.printf(ns, "sample %s => %v\n", a.displayList(ns, a.entries), outs)
	ns.dataSet.Data = append(ns.dataSet.Data, data)
}

//ExportSelectedAntSample .
func (ns *Nests) ExportSelectedAntSample() (int, error) {
	selected := ns.selected
	if selected <= 0 {
		return 0, nil
	}
	if ns.dataSet == nil {
		return 0, nil
	}
	if len(ns.dataSet.Data) == 0 {
		ns.printf("empty sample\n")
		return 0, nil
	}
	jsonString, _ := json.Marshal(ns.dataSet)
	err := ioutil.WriteFile("./tests/testant.json", jsonString, 0644)
	nn := len(ns.dataSet.Data)
	ns.dataSet.Data = make([]network.MlDataSample, 0, 0)
	return nn, err
}

//SetSleep .
func (ns *Nests) SetSleep(value int) {
	ns.waiter = value
	ns.printf("waiter set to: %d\n", ns.waiter)
}

//GetGlobalInfo .
func (ns *Nests) GetGlobalInfo() *GlobalInfo {
	return &GlobalInfo{
		Nests:        ns.nbs,
		Ndir:         outNb,
		Waiter:       ns.waiter,
		Xmin:         ns.xmin,
		Xmax:         ns.xmax,
		Ymin:         ns.ymin,
		Ymax:         ns.ymax,
		SelectedNest: ns.selectedNest,
		SelectedAnt:  ns.selected,
	}
}

//GetInfo .
func (ns *Nests) GetInfo() *Infos {

	ns.commitPeriodStats()
	global := &Info{
		Happiness:                   ns.happiness,
		Train:                       ns.statTrain.cumul / ns.totalNumber,
		Reinforce:                   ns.statReinforce.cumul / ns.totalNumber,
		Fade:                        ns.statFade.cumul / ns.totalNumber,
		Decision:                    ns.statDecision.cumul / ns.totalNumber,
		FromBeginningTrain:          ns.statTrain.scumul,
		Contact:                     ns.statContact.cumul,
		PeriodNetworkUpdated:        ns.statNetwork.cumul,
		FromBeginningNetworkUpdated: ns.statNetwork.scumul,
	}
	if ns.statDecision.cumul > 0 {
		global.PeriodGRate = float64(ns.statReinforce.cumul*100) / float64(ns.statDecision.cumul)
	}
	if ns.statDecision.scumul > 0 {
		global.FromBeginningGRate = float64(ns.statReinforce.scumul*100) / float64(ns.statDecision.scumul)
	}
	var selected = &Info{}
	aa := ns.getSelected()
	if aa != nil {
		aa.commitPeriodStats(ns)
		selected = &Info{
			Happiness:                   aa.happiness,
			Train:                       aa.statTrain.cumul,
			Reinforce:                   aa.statReinforce.cumul,
			Fade:                        aa.statFade.cumul,
			Decision:                    aa.statDecision.cumul,
			Contact:                     aa.statContact.cumul,
			PeriodNetworkUpdated:        aa.statNetwork.cumul,
			FromBeginningTrain:          aa.statTrain.scumul,
			FromBeginningNetworkUpdated: aa.statNetwork.scumul,
		}
		if ns.stopped {
			selected.Train = aa.statTrain.value
			selected.Reinforce = aa.statReinforce.value
			selected.Fade = aa.statFade.value
			selected.Decision = aa.statDecision.value
			selected.Contact = aa.statContact.value
			selected.PeriodNetworkUpdated = aa.statNetwork.value
		}
		if ns.stopped {
			if aa.statDecision.value > 0 {
				selected.PeriodGRate = float64(aa.statReinforce.value*100) / float64(aa.statDecision.value)
			}
		} else {
			if aa.statDecision.cumul > 0 {
				selected.PeriodGRate = float64(aa.statReinforce.cumul*100) / float64(aa.statDecision.cumul)
			}
		}
		if aa.statDecision.scumul > 0 {
			selected.FromBeginningGRate = float64(aa.statReinforce.scumul*100) / float64(aa.statDecision.scumul)
		}
	}
	ns.bestAnt.commitPeriodStats(ns)
	ns.worseAnt.commitPeriodStats(ns)

	return &Infos{
		Timer:                ns.timeRef,
		Speed:                ns.speed,
		Global:               global,
		Selected:             selected,
		BestID:               ns.bestAnt.id,
		BestNetworkStruct:    fmt.Sprintf("%v", ns.bestAnt.network.Getdef()),
		BestNetworkGRate:     ns.bestAnt.gRate,
		BestNetworkDirCount:  ns.bestAnt.dirCount,
		WorseID:              ns.worseAnt.id,
		WorseNetworkStruct:   fmt.Sprintf("%v", ns.worseAnt.network.Getdef()),
		WorseNetworkGRate:    ns.worseAnt.gRate,
		WorseNetworkDirCount: ns.worseAnt.dirCount,
		FromBeginningFoods:   ns.statFood.scumul,
		PeriodFoods:          ns.statFood.cumul,
	}
}

//GetNetwork .
func (ns *Nests) GetNetwork(nestID int, antID int) (*network.MLNetwork, error) {
	if nestID == 0 && antID == 0 {
		if ns.selected <= 0 {
			return nil, fmt.Errorf("no selected ant")
		}
		return ns.nests[ns.selectedNest-1].ants[ns.selected-1].network, nil
	}
	if nestID == -1 {
		return ns.bestAnt.network, nil
	}
	if nestID == -2 {
		return ns.worseAnt.network, nil
	}
	if nestID-1 < 0 || nestID-1 >= len(ns.nests) {
		return nil, fmt.Errorf("bad nest id: %d should be [1-%d]", nestID, len(ns.nests))
	}
	if antID-1 < 0 || antID-1 >= len(ns.nests[nestID-1].ants) {
		return nil, fmt.Errorf("bad ant id: %d should be [1-%d]", antID, len(ns.nests[nestID-1].ants))
	}
	return ns.nests[nestID-1].ants[antID-1].network, nil
}

// AddFoodGroup .
func (ns *Nests) AddFoodGroup(gx float64, gy float64) {
	fg := &FoodGroup{x: gx, y: gy}
	ns.foodGroups = append(ns.foodGroups, fg)
	for ii := 0; ii < 20; ii++ {
		food := newFood(fg)
		ns.foods = append(ns.foods, food)
	}
}

func (ns *Nests) addPheromone(x float64, y float64, id int) {
	var free *Pheromone
	for _, p := range ns.pheromones {
		if p.Level <= 0 {
			free = p
		}
		if (p.X-x)*(p.X-x)+(p.Y-y)*(p.Y-y) < ns.pheromoneGroup {
			p.Level = ns.pheromoneLevel
			return
		}
	}
	if free != nil {
		free.X = x
		free.Y = y
		free.Level = ns.pheromoneLevel
		free.id = id
		return
	}
	ns.pheromones = append(ns.pheromones, &Pheromone{X: x, Y: y, Level: ns.pheromoneLevel, id: id})
}

func (ns *Nests) fadePheromones() {
	ns.pheromoneFadeCounter--
	if ns.pheromoneFadeCounter < 0 {
		ns.pheromoneFadeCounter = ns.pheromoneFadeDelay
		for _, p := range ns.pheromones {
			p.Level--
		}
	}
}

//FoodRenew .
func (ns *Nests) FoodRenew(renew bool) {
	ns.foodRenew = renew
	if ns.foodRenew {
		for _, f := range ns.foods {
			fg := ns.foodGroups[rand.Int31n(int32(len(ns.foodGroups)))]
			fg.setPosition(f)
			f.carried = false
		}
	}
}
