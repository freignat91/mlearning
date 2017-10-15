package nests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"time"

	"github.com/freignat91/mlearning/network"
)

//Nests .
type Nests struct {
	waiter       int
	nests        []*Nest
	totalNumber  int
	xmin         float64
	xmax         float64
	ymin         float64
	ymax         float64
	stopped      bool
	timeRef      int64
	lastTimeRef  int64
	speed        int64
	selectedNest int
	selected     int
	averageRate  float64
	bestNest     *Nest
	ready        bool
	panicMode    bool
	happiness    float64
	dataSet      *network.MlDataSet
	foods        []*Food
	foodGroups   []*FoodGroup
	foodRenew    bool
	parameters   *Parameters
	//
	log           bool
	period        int64
	statTrain     *Stats
	statDecision  *Stats
	statReinforce *Stats
	statFade      *Stats
	statNetwork   *Stats
	statContact   *Stats
}

// GraphicData .
type GraphicData struct {
	Foods []*Food     `json:"foods"`
	Nests []*NestData `json:"nests"`
}

//GlobalInfo .
type GlobalInfo struct {
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
	Timer              int64       `json:"timer"`
	Speed              int64       `json:"speed"`
	Nests              []*NestInfo `json:"nests"`
	Selected           *Info       `json:"selected"`
	Global             *Info       `json:"global"`
	FromBeginningFoods int64       `json:"fromBeginningFoods"`
	PeriodFoods        int64       `json:"periodFoods"`
	FoodRenew          bool        `json:"foodRenew"`
	PanicMode          bool        `json:"panicMode"`
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
func NewNests(xmin float64, ymin float64, xmax float64, ymax float64, foodNb int, foodGroupNb int) (*Nests, error) {
	nests := &Nests{
		xmin:          xmin,
		ymin:          ymin,
		xmax:          xmax,
		ymax:          ymax,
		waiter:        1,
		parameters:    newParameters(),
		nests:         make([]*Nest, 2, 2),
		stopped:       true,
		selected:      0,
		selectedNest:  1,
		period:        10000,
		log:           false,
		foodRenew:     true,
		panicMode:     true,
		statTrain:     newStats(nil, nil),
		statDecision:  newStats(nil, nil),
		statReinforce: newStats(nil, nil),
		statFade:      newStats(nil, nil),
		statNetwork:   newStats(nil, nil),
		statContact:   newStats(nil, nil),
	}
	if err := nests.init(); err != nil {
		return nil, err
	}
	nests.ready = true
	return nests, nil
}

//IsReady .
func (ns *Nests) IsReady() bool {
	return ns.ready
}

func (ns *Nests) init() error {
	rand.Seed(time.Now().UTC().UnixNano())
	for ii, nest := range ns.nests {
		if nest == nil {
			nnest, err := newNest(ns, ii+1)
			if err != nil {
				return err
			}
			ns.nests[ii] = nnest
		} else {
			nest.init()
		}
	}
	ns.bestNest = ns.nests[0]
	ns.foods = make([]*Food, 0, 0)
	ns.foodGroups = make([]*FoodGroup, 0, 0)
	cc := 1.0
	for _, nest := range ns.nests {
		for ii := 0; ii < ns.parameters.initialFoodGroupNumberPerNest; ii++ {
			angle := 20 + rand.Float64()*50
			angle = angle * math.Pi / 180
			x := nest.x + cc*math.Cos(angle)*300
			y := nest.y + cc*math.Sin(angle)*300
			ns.AddFoodGroup(x, y)
		}
		cc = cc * (-1)
	}
	return nil
}

//GetGraphicData .
func (ns *Nests) GetGraphicData() *GraphicData {
	nests := make([]*NestData, len(ns.nests), len(ns.nests))
	for ii, nest := range ns.nests {
		nests[ii] = nest.getData()
	}
	return &GraphicData{
		Foods: ns.foods,
		Nests: nests,
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
			if ns.timeRef%10 == 0 {
				ns.verifRestart()
			}
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
	if ns.stopped {
		return
	}
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
	totalNb := 0
	for _, nest := range ns.nests {
		nest.nextTime(ns)
		totalNb += nest.workerNb + nest.soldierNb
		h += nest.happiness
		if nest.bestWorker.gRate > ns.bestNest.bestWorker.gRate {
			ns.bestNest = nest
		}
	}
	ns.totalNumber = totalNb
	ns.averageRate = float64(ns.statReinforce.cumul) * 100.0 / float64(ns.statDecision.cumul)
	ns.happiness = h / float64(len(ns.nests))
}

func (ns *Nests) printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (ns *Nests) getSelected() *Ant {
	if ns.selectedNest <= 0 || ns.selectedNest > len(ns.nests) {
		ns.selected = 0
		return nil
	}
	selectedNest := ns.nests[ns.selectedNest-1]
	if ns.selected <= 0 || ns.selected > len(selectedNest.ants) {
		ns.selected = 0
		return nil
	}
	return selectedNest.ants[ns.selected-1]
}

func (ns *Nests) commitPeriodStats() {
	if ns.stopped {
		return
	}
	ns.speed = ns.timeRef - ns.lastTimeRef
	ns.lastTimeRef = ns.timeRef
	for _, nest := range ns.nests {
		nest.setBestAnts()
		nest.statTrain.push()
		nest.statDecision.push()
		nest.statReinforce.push()
		nest.statFade.push()
		nest.statNetwork.push()
		nest.statContact.push()
	}
	ns.statTrain.push()
	ns.statDecision.push()
	ns.statReinforce.push()
	ns.statFade.push()
	ns.statNetwork.push()
	ns.statContact.push()
}

//SetSelected .
func (ns *Nests) SetSelected(selectedNest int, selected int) {
	ns.selected = selected
	ns.selectedNest = selectedNest
	ant := ns.getSelected()
	if ant != nil {
		fmt.Printf("selected: %d-%d type=%d life:%d\n", selectedNest, selected, ant.AntType, ant.Life)
		ns.dataSet = &network.MlDataSet{
			Name:   "ant",
			Layers: ant.network.Getdef(),
			Data:   make([]network.MlDataSample, 0, 0),
		}
		if ns.stopped {
			ant.commitPeriodStats(ns)
		}
	} else {
		ns.dataSet = nil
	}
}

func (ns *Nests) addSample(ins []float64, outs []float64) {
	if ns.dataSet == nil || len(ns.dataSet.Data) >= 10000 {
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
	if ns.totalNumber == 0 {
		ns.totalNumber = 1
	}
	global := &Info{
		Happiness:                   ns.happiness,
		Train:                       ns.statTrain.cumul / int64(ns.totalNumber),
		Reinforce:                   ns.statReinforce.cumul / int64(ns.totalNumber),
		Fade:                        ns.statFade.cumul / int64(ns.totalNumber),
		Decision:                    ns.statDecision.cumul / int64(ns.totalNumber),
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

	nests := make([]*NestInfo, len(ns.nests), len(ns.nests))
	for ii, nest := range ns.nests {
		nests[ii] = nest.getInfo()
	}
	return &Infos{
		Timer:     ns.timeRef,
		Speed:     ns.speed,
		Nests:     nests,
		Global:    global,
		Selected:  selected,
		FoodRenew: ns.foodRenew,
		PanicMode: ns.panicMode,
	}
}

//GetNetwork .
func (ns *Nests) GetNetwork(nestID int, antID int) (*network.MLNetwork, error) {
	if nestID == 0 && antID == 0 {
		ant := ns.getSelected()
		if ant == nil {
			return nil, fmt.Errorf("No selected ant")
		}
		return ant.network, nil
	}
	if ns.selectedNest <= 0 || ns.selectedNest > len(ns.nests) {
		return nil, fmt.Errorf("Invalid nest id")
	}
	selectedNest := ns.nests[ns.selectedNest-1]
	if nestID == -1 {
		return selectedNest.bestWorker.network, nil
	}
	if antID <= 0 || antID > len(selectedNest.ants) {
		return nil, fmt.Errorf("Invalid ant id")
	}
	return selectedNest.ants[ns.selected-1].network, nil
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

// RemoveFoodGroup .
func (ns *Nests) RemoveFoodGroup(gx float64, gy float64) {
	dist2m := 20.0 * 20.0
	var nfg *FoodGroup
	for _, fg := range ns.foodGroups {
		dist2 := (fg.x-gx)*(fg.x-gx) + (fg.y-gy)*(fg.y-gy)
		if dist2 < dist2m {
			dist2m = dist2
			nfg = fg
		}
	}
	if nfg != nil {
		//
	}
}

//FoodRenew .
func (ns *Nests) FoodRenew(renew bool) {
	if len(ns.foodGroups) == 0 {
		return
	}
	ns.foodRenew = renew
	if ns.foodRenew {
		for _, f := range ns.foods {
			f.renew()
			f.carried = false
		}
	}
}

//ClearFoodGroup .
func (ns *Nests) ClearFoodGroup() {
	ns.foodGroups = make([]*FoodGroup, 0, 0)
	ns.foods = make([]*Food, 0, 0)
}

//SetPanicMode .
func (ns *Nests) SetPanicMode(mode bool) {
	ns.panicMode = mode
}

func (ns *Nests) verifRestart() {
	if len(ns.nests[0].ants) == 0 {
		ns.nests[1].success++
		fmt.Printf("Nest red win worker=%d soldier=%d\n", ns.nests[1].workerNb, ns.nests[1].soldierNb)
		ns.init()
		return
	}
	if len(ns.nests[1].ants) == 0 {
		ns.nests[0].success++
		fmt.Printf("Nest blue win worker=%d soldier=%d\n", ns.nests[1].workerNb, ns.nests[1].soldierNb)
		ns.init()
		return
	}
}
