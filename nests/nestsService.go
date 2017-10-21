package nests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/freignat91/mlearning/network"
)

// GraphicData .
type GraphicData struct {
	Foods []*Food     `json:"foods"`
	Nests []*NestData `json:"nests"`
}

//GlobalInfo .
type GlobalInfo struct {
	Xmin         float64           `json:"xmin"`
	Xmax         float64           `json:"xmax"`
	Ymin         float64           `json:"ymin"`
	Ymax         float64           `json:"ymax"`
	Ndir         int               `json:"ndir"`
	Waiter       int               `json:"waiter"`
	SelectedNest int               `json:"selectedNest"`
	SelectedAnt  int               `json:"selectedAnt"`
	Nests        []*NestGlobalInfo `json:"nests"`
}

//Infos .
type Infos struct {
	Timer        int64         `json:"timer"`
	Speed        int64         `json:"speed"`
	Nests        []*NestInfo   `json:"nests"`
	FoodRenew    bool          `json:"foodRenew"`
	PanicMode    bool          `json:"panicMode"`
	SelectedInfo *SelectedInfo `json:"selectedInfo"`
}

//SelectedInfo .
type SelectedInfo struct {
	ID        int     `json:"id"`
	NestID    int     `json:"nestId"`
	Mode      string  `json:"mode"`
	Life      int     `json:"life"`
	Decision  int64   `json:"decision"`
	Reinforce int64   `json:"reinforce"`
	Fade      int64   `json:"fade"`
	DirCount  int     `json:"dirCount"`
	GRate     float64 `json:"gRate"`
}

//GetGlobalInfo .
func (ns *Nests) GetGlobalInfo() *GlobalInfo {
	ret := &GlobalInfo{
		Ndir:   outNb,
		Waiter: waiter,
		Xmin:   ns.xmin,
		Xmax:   ns.xmax,
		Ymin:   ns.ymin,
		Ymax:   ns.ymax,
	}
	nests := make([]*NestGlobalInfo, ns.nbNests, ns.nbNests)
	for ii, nest := range ns.nests {
		nests[ii] = &NestGlobalInfo{
			X: nest.x,
			Y: nest.y,
		}
	}
	ret.Nests = nests
	if ns.selected != nil {
		ret.SelectedNest = ns.selected.nest.id
		ret.SelectedAnt = ns.selected.ID
	}
	return ret
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

//SetSelected .
func (ns *Nests) SetSelected(selectedNest int, selected int, mode string) {
	ant := ns.getSelected(selectedNest, selected, mode)
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

//ExportSelectedAntSample .
func (ns *Nests) ExportSelectedAntSample() (int, error) {
	selected := ns.selected
	if selected == nil {
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
	waiter = value
	ns.printf("waiter set to: %d\n", waiter)
}

//IsReady .
func (ns *Nests) IsReady() bool {
	return ns.ready
}

//IsStarted .
func (ns *Nests) IsStarted() bool {
	return !ns.stopped
}

//GetSelectedNetwork .
func (ns *Nests) GetSelectedNetwork() (*network.MLNetwork, error) {
	if ns.selected != nil {
		return ns.selected.network, nil
	}
	return nil, fmt.Errorf("ant not selected")
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
	foodRenew = renew
	if foodRenew {
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
	panicMode = mode
}

// TrainSoluce .
func (ns *Nests) TrainSoluce(nb int) []string {
	if ns.selected == nil {
		return []string{"No ant selected\n"}
	}
	ns.selected.trainSoluce(ns, nb)
	return ns.selected.test(ns)
}

// Test .
func (ns *Nests) Test() []string {
	if ns.selected == nil {
		return []string{"No ant selected\n"}
	}
	return ns.selected.test(ns)
}

// LogsToggle .
func (ns *Nests) LogsToggle() []string {
	ns.log = !ns.log
	if ns.log {
		return []string{"Server logs on\n"}
	}
	return []string{"Server logs off\n"}
}

//GetInfo .
func (ns *Nests) GetInfo() *Infos {
	aa := ns.selected
	var selectedInfo *SelectedInfo
	if aa != nil {
		selectedInfo = &SelectedInfo{
			ID:        aa.ID,
			NestID:    aa.nest.id,
			Mode:      aa.getModeToString(),
			Life:      aa.Life,
			Decision:  aa.statDecision.cumul,
			Reinforce: aa.statReinforce.cumul,
			Fade:      aa.statFade.cumul,
			DirCount:  aa.dirCount,
			GRate:     aa.gRate,
		}
	}
	nests := make([]*NestInfo, len(ns.nests), len(ns.nests))
	for ii, nest := range ns.nests {
		nests[ii] = nest.getInfo()
	}
	diff := time.Now().Sub(ns.lastUpdateTime).Seconds()
	ns.lastUpdateTime = time.Now()
	ns.speed = int64(float64(ns.timeRef-ns.lastTimeRef) / diff)
	ns.lastTimeRef = ns.timeRef
	ret := &Infos{
		Timer:        ns.timeRef,
		Speed:        ns.speed,
		Nests:        nests,
		FoodRenew:    foodRenew,
		PanicMode:    panicMode,
		SelectedInfo: selectedInfo,
	}
	return ret
}
