package nests

import (
	"fmt"
	"math/rand"
)

//Nest .
type Nest struct {
	ns                   *Nests
	id                   int
	x                    float64
	y                    float64
	antIDCounter         int
	ants                 []*Ant
	statTrain            *Stats
	statDecision         *Stats
	statReinforce        *Stats
	statFade             *Stats
	statNetwork          *Stats
	statContact          *Stats
	bestWorker           *Ant
	bestSoldier          *Ant
	happiness            float64
	happinessTmp         float64
	parameters           *Parameters
	pheromones           []*Pheromone
	pheromoneFadeCounter int
	ressource            int
	workerNb             int
	soldierNb            int
	life                 int
	lifeTmp              int
	bestDirCount         int
	bestGRate            float64
	success              int
}

//NestData .
type NestData struct {
	Ants       []*Ant       `json:"ants"`
	Pheromones []*Pheromone `json:"pheromones"`
}

//NestInfo .
type NestInfo struct {
	NestID              int     `json:"id"`
	Worker              int     `json:"worker"`
	Soldier             int     `json:"soldier"`
	Ressource           int     `json:"ressource"`
	Life                int     `json:"life"`
	BestNetworkStruct   string  `json:"bestNetworkStruct"`
	BestNetworkDirCount int     `json:"bestNetworkDirCount"`
	BestNetworkGRate    float64 `json:"bestNetworkGRate"`
	Success             int     `json:"success"`
}

func newNest(ns *Nests, id int) (*Nest, error) {
	param := newParameters()
	nest := &Nest{
		id:            id,
		ns:            ns,
		parameters:    param,
		statTrain:     newStats(nil, nil),
		statDecision:  newStats(nil, nil),
		statReinforce: newStats(nil, nil),
		statFade:      newStats(nil, nil),
		statNetwork:   newStats(nil, nil),
		statContact:   newStats(nil, nil),
	}
	if id == 1 {
		nest.x = ns.xmin + 25
		nest.y = ns.ymin + 25
	} else if id == 2 {
		nest.x = ns.xmax - 25
		nest.y = ns.ymax - 25
	} else if id == 3 {
		nest.x = ns.xmin + 25
		nest.y = ns.ymax - 25
	} else {
		nest.x = ns.xmax - 25
		nest.y = ns.ymin + 25
	}
	nest.init()
	return nest, nil
}

func (n *Nest) getData() *NestData {
	return &NestData{
		Ants:       n.ants,
		Pheromones: n.pheromones,
	}
}

func (n *Nest) init() {
	n.workerNb = 0
	n.soldierNb = 0
	n.ants = make([]*Ant, 0, 0)
	n.pheromones = make([]*Pheromone, 0, 0)
	n.ressource = n.parameters.workerInitNb
	for ii := 0; ii < n.parameters.workerInitNb; ii++ {
		n.addWorker(n.ns, n.x+30-rand.Float64()*60, n.y+30-rand.Float64()*60, rand.Intn(outNb))
	}
	n.ressource = n.parameters.soldierInitNb
	for ii := 0; ii < n.parameters.soldierInitNb; ii++ {
		n.addSoldier(n.ns, n.x+30-rand.Float64()*60, n.y+30-rand.Float64()*60, 0, 0, rand.Intn(outNb))
	}
	n.bestWorker = newAntWorker(n.ns, n, 0, 0, 0)
	n.bestSoldier = newAntSoldier(n.ns, n, 0, 0, 0, 0, 0)
	n.ressource = n.parameters.nestInitialRessource
}

func (n *Nest) getInfo() *NestInfo {
	return &NestInfo{
		NestID:              n.id,
		Worker:              n.workerNb,
		Soldier:             n.soldierNb,
		Ressource:           n.ressource,
		Life:                n.life,
		BestNetworkStruct:   fmt.Sprintf("%v", n.bestWorker.network.Getdef()),
		BestNetworkDirCount: n.bestDirCount,
		BestNetworkGRate:    n.bestGRate,
		Success:             n.success,
	}
}

func (n *Nest) nextTime(ns *Nests) {
	n.happinessTmp = 0
	n.lifeTmp = 0
	n.fadePheromones()
	if n.parameters.nestAntRenewDelay > 0 && ns.timeRef%n.parameters.nestAntRenewDelay == 0 {
		n.ressource++
		n.addWorker(ns, n.x+30-rand.Float64()*60, n.y+30-rand.Float64()*60, rand.Intn(outNb))
	}
	for _, ant := range n.ants {
		n.lifeTmp += ant.Life
		n.happinessTmp += ant.happiness
		ant.nextTime(ns)
	}
	n.life = n.lifeTmp
	n.happiness = 0
	if len(n.ants) > 0 {
		n.happiness = n.happinessTmp / float64(len(n.ants))
	}
	n.removeDeadAnts()
}

func (n *Nest) setBestAnts() {
	for _, ant := range n.ants {
		ant.commitPeriodStats(n.ns)
		if ant.gRate < 100 {
			if ant.dirCount > n.bestDirCount {
				n.setBestAnt(ant)
			} else if ant.dirCount == n.bestDirCount && ant.gRate > n.bestGRate {
				n.setBestAnt(ant)
			}
		}
	}
}

func (n *Nest) setBestAnt(ant *Ant) {
	n.bestDirCount = ant.dirCount
	n.bestGRate = ant.gRate
	if ant.AntType == 0 {
		n.bestWorker = ant
	} else {
		n.bestSoldier = ant
	}
}

func (n *Nest) addWorker(ns *Nests, x float64, y float64, direction int) {
	if n.ressource > 0 && n.workerNb < n.parameters.maxWorkerAnt && n.workerNb+n.soldierNb < n.parameters.maxAnt {
		ant := newAntWorker(ns, n, x, y, direction)
		if ant != nil {
			n.ants = append(n.ants, ant)
			n.ressource--
		}
	}
}

func (n *Nest) addSoldier(ns *Nests, x float64, y float64, dx float64, dy float64, direction int) {
	if n.ressource >= n.parameters.soldierRessourceCost && n.workerNb+n.soldierNb < n.parameters.maxAnt {
		ant := newAntSoldier(ns, n, x, y, dx, dy, direction)
		if ant != nil {
			n.ants = append(n.ants, ant)
			n.ressource -= n.parameters.soldierRessourceCost
		}
	}
}

func (n *Nest) removeDeadAnts() {
	workerNb := 0
	soldierNb := 0
	ants := make([]*Ant, 0, 0)
	for _, ant := range n.ants {
		if ant.Life > 0 {
			ants = append(ants, ant)
			if ant.AntType == 0 {
				workerNb++
			} else {
				soldierNb++
			}
		} else {
			ant.dropFood(n.ns)
		}
	}
	n.workerNb = workerNb
	n.soldierNb = soldierNb
	n.ants = ants
}

func (n *Nest) addPheromone(x float64, y float64, id int) {
	var free *Pheromone
	for _, p := range n.pheromones {
		if p.Level <= 0 {
			free = p
		}
		if (p.X-x)*(p.X-x)+(p.Y-y)*(p.Y-y) < n.parameters.pheromoneGroup {
			p.Level = n.parameters.pheromoneLevel
			return
		}
	}
	if free != nil {
		free.X = x
		free.Y = y
		free.Level = n.parameters.pheromoneLevel
		free.id = id
		return
	}
	phe := &Pheromone{X: x, Y: y, Level: n.parameters.pheromoneLevel, id: id}
	n.pheromones = append(n.pheromones, phe)
}

func (n *Nest) fadePheromones() {
	n.pheromoneFadeCounter--
	if n.pheromoneFadeCounter < 0 {
		n.pheromoneFadeCounter = n.parameters.pheromoneFadeDelay
		for _, p := range n.pheromones {
			p.Level--
		}
	}
}
