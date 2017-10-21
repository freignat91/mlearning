package nests

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/freignat91/mlearning/network"
)

var waiter = 1
var foodRenew = true
var panicMode = true

//Nests .
type Nests struct {
	nests       []*Nest
	xmin        float64
	xmax        float64
	ymin        float64
	ymax        float64
	stopped     bool
	timeRef     int64
	lastTimeRef int64
	speed       int64
	selected    *Ant
	ready       bool
	dataSet     *network.MlDataSet
	foods       []*Food
	foodGroups  []*FoodGroup
	parameters  *Parameters
	nbNests     int
	//
	log            bool
	period         int64
	nextMutex      sync.RWMutex
	lastUpdateTime time.Time
}

//NewNests .
func NewNests(xmin float64, ymin float64, xmax float64, ymax float64, foodNb int, foodGroupNb int, nbNests int) (*Nests, error) {
	param := newParameters()
	nests := &Nests{
		xmin:           xmin,
		ymin:           ymin,
		xmax:           xmax,
		ymax:           ymax,
		parameters:     param,
		nbNests:        nbNests,
		nests:          make([]*Nest, nbNests, nbNests),
		stopped:        true,
		period:         10000,
		log:            false,
		lastUpdateTime: time.Now(),
	}
	if err := nests.init(); err != nil {
		return nil, err
	}
	nests.ready = true
	return nests, nil
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
	//
	//ns.selected = ns.nests[0].ants[10]
	ns.foods = make([]*Food, 0, 0)
	ns.foodGroups = make([]*FoodGroup, 0, 0)
	dx := []float64{1, -1, 1, -1}
	dy := []float64{1, -1, -1, 1}
	ll := ns.parameters.initial2FoodGroupLenght
	if ns.nbNests == 4 {
		ll = ns.parameters.initial4FoodGroupLenght
	}
	for ii, nest := range ns.nests {
		for jj := 0; jj < ns.parameters.initialFoodGroupNumberPerNest; jj++ {
			angle := 20 + rand.Float64()*50
			angle = angle * math.Pi / 180
			x := nest.x + dx[ii]*math.Cos(angle)*float64(ll)
			y := nest.y + dy[ii]*math.Sin(angle)*float64(ll)
			ns.AddFoodGroup(x, y)
		}
	}
	return nil
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
			if waiter > 0 {
				time.Sleep(time.Duration(waiter) * time.Millisecond)
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
	if ns.selected != nil && ns.selected.Life > 0 {
		ns.selected.commitPeriodStats(ns)
	}
	ns.stopped = true
}

//NextTime .
func (ns *Nests) NextTime() {
	//ns.printf("\033[2J\033[0;0H\n")
	ns.nextMutex.Lock()
	ns.timeRef++
	update := false
	if ns.timeRef%ns.parameters.updateTickNumber == 0 {
		update = true
	}
	for _, nest := range ns.nests {
		nest.nextTime(ns, update)
	}
	if ns.selected != nil && ns.selected.Life <= 0 {
		ns.selected = nil
	}
	ns.nextMutex.Unlock()
}

func (ns *Nests) printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (ns *Nests) getSelected(nestID int, antID int, mode string) *Ant {
	if strings.ToLower(mode) == "selected" {
		return ns.selected
	}
	if nestID <= 0 || nestID > len(ns.nests) {
		ns.selected = nil
		return nil
	}
	selectedNest := ns.nests[nestID-1]
	if strings.ToLower(mode) == "bestworker" {
		ns.selected = selectedNest.bestWorker
		return ns.selected
	}
	if strings.ToLower(mode) == "bestsoldier" {
		ns.selected = selectedNest.bestSoldier
		return ns.selected
	}
	ns.selected = nil
	for _, ant := range selectedNest.ants {
		if ant.ID == antID {
			ns.selected = ant
			break
		}
	}
	return ns.selected
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

func (ns *Nests) verifRestart() {
	end := 0
	for _, nest := range ns.nests {
		if len(nest.ants) == 0 {
			end++
		}
	}
	if end == ns.nbNests-1 {
		for ii, nest := range ns.nests {
			if len(nest.ants) > 0 {
				nest.success++
				fmt.Printf("Nest %d win worker=%d soldier=%d\n", ii, ns.nests[ii].workerNb, ns.nests[ii].soldierNb)
			}
		}
		ns.init()
	}
}
