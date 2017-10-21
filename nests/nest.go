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
	statDecision         *Stats
	statGoodDecision     *Stats
	statReinforce        *Stats
	statFade             *Stats
	bestWorker           *Ant
	bestSoldier          *Ant
	parameters           *Parameters
	pheromones           []*Pheromone
	pheromoneFadeCounter int
	ressource            int
	workerNb             int
	soldierNb            int
	bestWorkerDirCount   int
	bestWorkerGRate      float64
	bestSoldierDirCount  int
	bestSoldierGRate     float64
	success              int
	averageDirCount      float64
	averageGRate         float64
}

//NestData .
type NestData struct {
	Ants       []*Ant       `json:"ants"`
	Pheromones []*Pheromone `json:"pheromones"`
}

//NestInfo .
type NestInfo struct {
	NestID                     int     `json:"id"`
	Worker                     int     `json:"worker"`
	Soldier                    int     `json:"soldier"`
	Ressource                  int     `json:"ressource"`
	BestWorkerNetworkStruct    string  `json:"bestWorkerNetworkStruct"`
	BestWorkerNetworkDirCount  int     `json:"bestWorkerNetworkDirCount"`
	BestWorkerNetworkGRate     float64 `json:"bestWorkerNetworkGRate"`
	BestSoldierNetworkStruct   string  `json:"bestSoldierNetworkStruct"`
	BestSoldierNetworkDirCount int     `json:"bestSoldierNetworkDirCount"`
	BestSoldierNetworkGRate    float64 `json:"bestSoldierNetworkGRate"`
	Success                    int     `json:"success"`
	Decision                   float64 `json:"decision"`
	Reinforce                  float64 `json:"reinforce"`
	Fade                       float64 `json:"fade"`
	DirCount                   float64 `json:"dirCount"`
	GRate                      float64 `json:"gRate"`
}

//NestGlobalInfo .
type NestGlobalInfo struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func newNest(ns *Nests, id int) (*Nest, error) {
	param := newParameters()
	nest := &Nest{
		id:               id,
		ns:               ns,
		parameters:       param,
		statDecision:     newStats(nil),
		statGoodDecision: newStats(nil),
		statReinforce:    newStats(nil),
		statFade:         newStats(nil),
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
	n.ressource = n.parameters.soldierInitNb * n.parameters.soldierRessourceCost
	for ii := 0; ii < n.parameters.soldierInitNb; ii++ {
		n.addSoldier(n.ns, n.x+30-rand.Float64()*60, n.y+30-rand.Float64()*60, 0, 0, rand.Intn(outNb))
	}
	if n.bestWorker == nil {
		n.bestWorker = newAntWorker(n.ns, n, 0, 0, 0)
	}
	if n.bestSoldier == nil {
		n.bestSoldier = newAntSoldier(n.ns, n, 0, 0, 0, 0, 0)
	}
	n.ressource = n.parameters.nestInitialRessource
}

func (n *Nest) commitStats() {
	n.statDecision.push()
	n.statGoodDecision.push()
	n.statReinforce.push()
	n.statFade.push()
}

func (n *Nest) getInfo() *NestInfo {
	nb := float64(n.workerNb + n.soldierNb)
	ret := &NestInfo{
		NestID:                     n.id,
		Worker:                     n.workerNb,
		Soldier:                    n.soldierNb,
		Ressource:                  n.ressource,
		BestWorkerNetworkStruct:    fmt.Sprintf("%v", n.bestWorker.network.Getdef()),
		BestWorkerNetworkDirCount:  n.bestWorker.dirCount,
		BestWorkerNetworkGRate:     n.bestWorker.gRate,
		BestSoldierNetworkStruct:   fmt.Sprintf("%v", n.bestSoldier.network.Getdef()),
		BestSoldierNetworkDirCount: n.bestSoldier.dirCount,
		BestSoldierNetworkGRate:    n.bestSoldier.gRate,
		Success:                    n.success,
		DirCount:                   n.averageDirCount,
		GRate:                      n.averageGRate,
	}
	if nb > 0 {
		ret.Decision = float64(n.statDecision.cumul) / nb
		ret.Reinforce = float64(n.statReinforce.cumul) / nb
		ret.Fade = float64(n.statFade.cumul) / nb
	}
	return ret
}

func (n *Nest) nextTime(ns *Nests, update bool) {
	dirCountTmp := 0
	gRateTmp := 0.0
	n.fadePheromones()
	if n.parameters.nestAntRenewDelay > 0 && ns.timeRef%n.parameters.nestAntRenewDelay == 0 {
		n.ressource++
		n.addWorker(ns, n.x+30-rand.Float64()*60, n.y+30-rand.Float64()*60, rand.Intn(outNb))
	}
	if update {
		n.setBestAnts()
	}
	for _, ant := range n.ants {
		ant.nextTime(ns, update)
		dirCountTmp += ant.dirCount
		gRateTmp += ant.gRate
	}
	n.averageDirCount = 0
	n.averageGRate = 0
	if n.workerNb+n.soldierNb > 0 {
		n.averageDirCount = float64(dirCountTmp) / float64(n.workerNb+n.soldierNb)
		n.averageGRate = gRateTmp / float64(n.workerNb+n.soldierNb)
	}

	if update {
		n.removeDeadAnts()
		n.removePheromones()
		n.commitStats()
	}
}

func (n *Nest) setBestAnts() {
	if len(n.ants) == 0 {
		return
	}
	//n.bestWorkerDirCount = 0
	//n.bestWorkerGRate = 0
	//n.bestSoldierDirCount = 0
	//n.bestSoldierGRate = 0
	for _, ant := range n.ants {
		ant.commitPeriodStats(n.ns)
	}
	for _, ant := range n.ants {
		if ant.AntType == 0 {
			if ant.dirCount > n.bestWorkerDirCount {
				n.setBestAnt(ant)
			} else if ant.dirCount == n.bestWorkerDirCount && ant.gRate >= n.bestWorkerGRate {
				n.setBestAnt(ant)
			}
		} else {
			if ant.dirCount > n.bestSoldierDirCount {
				n.setBestAnt(ant)
			} else if ant.dirCount == n.bestSoldierDirCount && ant.gRate >= n.bestSoldierGRate {
				n.setBestAnt(ant)
			}
		}
	}
}

func (n *Nest) setBestAnt(ant *Ant) {
	if ant.AntType == 0 {
		n.bestWorkerDirCount = ant.dirCount
		n.bestWorkerGRate = ant.gRate
		n.bestWorker = ant
	} else {
		n.bestSoldierDirCount = ant.dirCount
		n.bestSoldierGRate = ant.gRate
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

func (n *Nest) removePheromones() {
	phes := make([]*Pheromone, 0, 0)
	for _, phe := range n.pheromones {
		if phe.Level > 0 {
			phes = append(phes, phe)
		}
	}
	n.pheromones = phes
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
