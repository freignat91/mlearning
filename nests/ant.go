package nests

import (
	"math"
	"math/rand"

	"github.com/freignat91/mlearning/network"
)

var visionNb = 8
var inNb = 8 * 4 //8 for ants, 8 for foods, 8 for pheromones, 8 for hostile
var outNb = 8

// Ant .
type Ant struct {
	ID        int     `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Direction int     `json:"direction"`
	Contact   bool    `json:"contact"`
	Fight     bool    `json:"fight"`
	AntType   int     `json:"type"` //worker=0, soldier=1
	Life      int     `json:"life"`
	//
	nest               *Nest
	happiness          float64
	lastHappiness      float64
	regularSpeed       float64
	maxSpeed           float64
	speed              float64
	vision             float64
	dx                 float64
	dy                 float64
	soldierInitCounter int
	network            *network.MLNetwork
	entries            []float64
	lastEntries        []float64
	outs               []float64
	lifeTime           int64
	lastDecision       int
	panic              bool
	statTrain          *Stats
	statDecision       *Stats
	statReinforce      *Stats
	statFade           *Stats
	statNetwork        *Stats
	statContact        *Stats
	gRate              float64
	dirMap             map[int]int
	dirCount           int
	food               *Food
	pheromoneDelay     int
	pheromoneCount     int
	lastPheromone      int
	lastPheromoneCount int
	entryMode          int
	lastEntryMode      int
	timeWithoutHostile int
	pursue             *Ant
	//updateMutex        sync.RWMutex
	//
	decTmp []int
}

func newAnt(ns *Nests, n *Nest, antType int) *Ant {
	n.antIDCounter++
	ant := &Ant{
		ID:            n.antIDCounter,
		X:             n.x + 20.0 - rand.Float64()*40,
		Y:             n.y + 20.0 - rand.Float64()*40,
		AntType:       antType,
		nest:          n,
		vision:        30,
		Direction:     int(rand.Int31n(int32(outNb))),
		entries:       make([]float64, inNb, inNb),
		lastEntries:   make([]float64, inNb, inNb),
		outs:          make([]float64, outNb, outNb),
		lastDecision:  -1,
		statTrain:     newStats(ns.statTrain, n.statTrain),
		statDecision:  newStats(ns.statDecision, n.statDecision),
		statReinforce: newStats(ns.statReinforce, n.statReinforce),
		statFade:      newStats(ns.statFade, n.statFade),
		statNetwork:   newStats(ns.statNetwork, n.statNetwork),
		statContact:   newStats(ns.statContact, n.statContact),
		dirMap:        make(map[int]int),
		decTmp:        make([]int, outNb/2, outNb/2),
	}
	ant.setNetwork(ns)
	return ant
}

func newAntWorker(ns *Nests, n *Nest, x float64, y float64, direction int) *Ant {
	ant := newAnt(ns, n, 0)
	ant.X = x
	ant.Y = y
	ant.Life = n.parameters.workerLife
	ant.regularSpeed = n.parameters.workerMinSpeed + rand.Float64()*.1
	ant.maxSpeed = n.parameters.workerMaxSpeed + rand.Float64()*.1
	ant.speed = ant.regularSpeed
	ant.Direction = direction
	return ant
}

func newAntSoldier(ns *Nests, n *Nest, x float64, y float64, dx float64, dy float64, direction int) *Ant {
	ant := newAnt(ns, n, 1)
	ant.X = x
	ant.Y = y
	ant.dx = dx
	ant.dy = dy
	if dx != 0 && dy != 0 {
		ant.soldierInitCounter = n.parameters.soldierInitCounter
	}
	ant.Life = n.parameters.soldierLife
	ant.regularSpeed = n.parameters.soldierMinSpeed + rand.Float64()*.1
	ant.maxSpeed = n.parameters.soldierMaxSpeed + rand.Float64()*.1
	ant.speed = ant.regularSpeed
	ant.Direction = direction
	return ant
}

func (a *Ant) decrLife(ns *Nests, val int) {
	if a.Life > 0 {
		a.Life -= val
		if a.Life <= 0 {
			a.Life = 0
		}
	}
}

func (a *Ant) setNetwork(ns *Nests) {
	if a.AntType == 0 {
		if a.nest.bestWorker != nil && rand.Float64() < 0.8 {
			net, err := a.nest.bestWorker.network.Copy()
			if err != nil {
				a.network = net
				return
			}
		}
	} else {
		if a.nest.bestSoldier != nil && rand.Float64() < 0.8 {
			net, err := a.nest.bestSoldier.network.Copy()
			if err != nil {
				a.network = net
				return
			}
		}
	}
	var defnet []int
	if true { //rand.Float64() < 0.5 {
		defnet = make([]int, 3, 3)
		defnet[0] = inNb
		defnet[1] = int(5 + rand.Int31n(45))
		defnet[2] = outNb
	} else {
		defnet = make([]int, 4, 4)
		defnet[0] = inNb
		defnet[1] = int(5 + rand.Int31n(45))
		defnet[2] = int(5 + rand.Int31n(25))
		defnet[3] = outNb
	}
	net, _ := network.NewNetwork(defnet)
	a.network = net
}

func (a *Ant) nextTime(ns *Nests) {
	if !a.tickLife(ns) {
		return
	}
	a.displayInfo(ns)
	a.moveOnOut(ns)
	if a.panic {
		return
	}
	a.updateEntries(ns)
	a.computeHappiness(ns)
	a.printf(ns, "happiness=%.3f entries: %s\n", a.happiness, a.displayList(ns, a.entries, "%.2f"))
	if ns.log && ns.selected == a.ID {
		a.printf(ns, "lastDecision:%d entrie: %s hapiness=%.5f delta=%f\n", a.lastDecision, a.displayList(ns, a.entries, "%.3f"), a.happiness, a.happiness-a.lastHappiness)
	}
	if a.happiness == a.lastHappiness && a.happiness >= 0 {
		a.printf(ns, "decision: no need\n")
	} else if a.happiness < a.lastHappiness { //,= ?
		a.fadeLastDecision(ns)
		if a.decide(ns) {
			if ns.log && ns.selected == a.ID {
				a.printf(ns, "decision using network: %d outs: %s\n", a.Direction, a.displayList(ns, a.outs, "%.3f"))
			}
			a.statDecision.incr()
		} else {
			a.Direction = int(rand.Int31n(int32(outNb)))
			a.lastDecision = -1
			a.printf(ns, "decision random: %d\n", a.Direction)
		}
	} else {
		a.printf(ns, "decision: no need\n")
		if a.train(ns) {
			a.statTrain.incr()
			if a.lastDecision != -1 {
				a.printf(ns, "Decision %d reinforced\n", a.lastDecision)
				a.statReinforce.incr()
				a.lastDecision = -1
			}
		}
	}
	a.update(ns)
}

func (a *Ant) tickLife(ns *Nests) bool {
	if a.Life < 0 {
		return false
	}
	if a.AntType == 0 {
		if ns.timeRef%a.nest.parameters.workerLifeDecrPeriod == 0 {
			a.decrLife(ns, 1)
			if a.Life <= 0 {
				return false
			}
		}
		return true
	}
	if ns.timeRef%a.nest.parameters.soldierLifeDecrPeriod == 0 {
		a.decrLife(ns, 1)
		if a.Life <= 0 {
			return false
		}
	}
	return true
}

func (a *Ant) decide(ns *Nests) bool {
	if rand.Float64() < 0.1 {
		//return false
	}
	if a.food != nil {
		return false
	}
	ins, ok := a.preparedEntries(a.entries)
	if !ok {
		a.printf(ns, "bad entries: %s\n", a.displayList(ns, a.entries, "%.3f"))
		return false
	}
	a.outs = a.network.Propagate(ins, true)
	if ns.log {
		a.printf(ns, "Compute decision, propagation: %s\n", a.displayList(ns, a.outs, "%.3f"))
	}
	direction := 0
	max := 0.0
	for ii, out := range a.outs {
		if out > max {
			max = out
			direction = ii
		}
	}
	ref := a.dirMap[direction]
	dp := a.getDirIndex(direction + 1)
	if a.dirMap[dp] < ref/2 {
		direction = dp
	}
	dp = a.getDirIndex(direction - 1)
	if a.dirMap[dp] <= ref/2 {
		direction = dp
	}
	a.Direction = direction
	a.lastDecision = direction
	a.dirMap[direction]++
	return true
}

func (a *Ant) update(ns *Nests) {
	a.lifeTime++
	if a.lifeTime%20000 == 0 && a.dirCount > 0 {
		if false { //a.dirCount < 3 {
			a.setNetwork(ns)
			a.statNetwork.incr()
			a.printf(ns, "recreate new random network: %v\n", a.network.Getdef())
		} else {
			if a.AntType == 0 {
				if a.dirCount < a.nest.bestWorker.dirCount-a.nest.parameters.networkUpdateDirCountDiff || (a.dirCount <= a.nest.bestWorker.dirCount && a.gRate < a.nest.bestWorker.gRate-a.nest.parameters.networkUpdateGRateDiff) {
					net, err := a.nest.bestWorker.network.Copy()
					if err == nil {
						a.network = net
						a.lastDecision = -1
						a.statNetwork.incr()
						a.printf(ns, "update network with the best worker: %v\n", a.network.Getdef())
					}
				}
			} else {
				if a.dirCount < a.nest.bestSoldier.dirCount-a.nest.parameters.networkUpdateDirCountDiff || (a.dirCount <= a.nest.bestSoldier.dirCount && a.gRate < a.nest.bestSoldier.gRate-a.nest.parameters.networkUpdateGRateDiff) {
					net, err := a.nest.bestSoldier.network.Copy()
					if err == nil {
						a.network = net
						a.lastDecision = -1
						a.statNetwork.incr()
						a.printf(ns, "update network with the best soldier: %v\n", a.network.Getdef())
					}
				}
			}
		}
	}
}

func (a *Ant) displayInfo(ns *Nests) {
	if ns.timeRef%10000 == 0 {
		a.dirMap = make(map[int]int)
		a.dirCount = a.network.ComputeDistinctOut()
	}
	if ns.log && a.ID == ns.selected {
		ggRate := float64(a.statReinforce.scumul) * 100.0 / float64(a.statDecision.scumul)
		a.printf(ns, "[%d] totTrain: %d train:%d reinforce:%d fade:%d decision:%d period:good=%.2f%%) global:good=%.2f%%)\n", a.ID, a.statTrain.scumul, a.statTrain.cumul, a.statReinforce.cumul, a.statFade.cumul, a.statDecision.cumul, a.gRate, ggRate)
		a.printf(ns, "network=%v hapiness=%.5f move: %d\n", a.network.Getdef(), a.happiness, a.Direction)
	}
}

func (a *Ant) commitPeriodStats(ns *Nests) {
	if ns.stopped {
		return
	}
	if a.statDecision.value != 0 {
		a.statTrain.push()
		a.statDecision.push()
		a.statReinforce.push()
		a.statFade.push()
		a.statNetwork.push()
		a.statContact.push()
		a.gRate = float64(a.statReinforce.cumul) * 100.0 / float64(a.statDecision.cumul)
	}
}

func (a *Ant) updateEntries(ns *Nests) {
	a.Contact = false
	a.Fight = false
	for ii := range a.entries {
		a.lastEntries[ii] = a.entries[ii]
		a.entries[ii] = 0
	}
	a.lastEntryMode = a.entryMode
	a.entryMode = 0
	if a.food != nil {
		return
	}

	if a.updateEntriesForHostileAnts(ns) {
		return
	}
	if a.AntType == 0 {
		if a.updateEntriesForFoods(ns) {
			return
		}
		if a.updateEntriesForPheromones(ns) {
			return
		}
	}
	a.updateEntriesForFriendAnts(ns)
	return
}

func (a *Ant) updateEntriesForFoods(ns *Nests) bool {
	dist2Max := a.vision * a.vision
	dist2m := dist2Max
	var foodMin *Food
	for _, food := range ns.foods {
		if !food.carried {
			dist2 := (food.X-a.X)*(food.X-a.X) + (food.Y-a.Y)*(food.Y-a.Y)
			if dist2 < dist2m {
				foodMin = food
				dist2m = dist2
			}
		}
	}
	a.printf(ns, "closest food: %+v\n", foodMin)
	if foodMin != nil {
		if dist2m < 4 {
			a.carryFood(foodMin)
			a.pheromoneCount = 0
			return true
		}
		ang := math.Atan2(foodMin.X-a.X, foodMin.Y-a.Y)
		if ang < 0 {
			ang = 2*math.Pi + ang
		}
		index := int(ang*float64(visionNb)/2.0/math.Pi + 0.000001)
		if index >= visionNb {
			index = index - visionNb
		}
		a.entries[visionNb+index] = ((dist2Max - dist2m) / dist2Max)
		return true
	}
	return false
}

func (a *Ant) updateEntriesForPheromones(ns *Nests) bool {
	//search pheromones
	minLevel := 100000
	dist2Max := a.vision * a.vision
	dist2m := dist2Max
	var pheMin *Pheromone
	for _, phe := range a.nest.pheromones {
		if phe.Level > 0 {
			dist2 := (phe.X-a.X)*(phe.X-a.X) + (phe.Y-a.Y)*(phe.Y-a.Y)
			if dist2 < dist2Max*1.5*1.5 && phe.id < minLevel {
				dist2m = dist2
				pheMin = phe
				minLevel = phe.id
			}
		}
	}
	if pheMin != nil {
		if a.lastPheromone == pheMin.id {
			a.lastPheromoneCount++
		} else {
			a.lastPheromoneCount = 0
		}
		a.lastPheromone = pheMin.id
		if a.lastPheromoneCount < 50 {
			ang := math.Atan2(pheMin.X-a.X, pheMin.Y-a.Y)
			if ang < 0 {
				ang = 2*math.Pi + ang
			}
			index := int(ang*float64(visionNb)/2.0/math.Pi + 0.000001)
			if index >= visionNb {
				index = index - visionNb
			}
			a.printf(ns, "better phe: %+v angle=%.4f index=%d level=%0.3f id=%d\n", pheMin, ang*180/math.Pi, index, pheMin.Level, pheMin.id)
			a.entries[visionNb*2+index] = ((dist2Max - dist2m) / dist2Max)
		}
		return true
	}
	return false
}

func (a *Ant) updateEntriesForHostileAnts(ns *Nests) bool {
	if a.AntType == 0 && !ns.panicMode {
		return false
	}
	dist2Max := a.vision * a.vision
	if a.AntType == 1 && a.pursue != nil {
		if a.pursue.Life > 0 {
			dist2m := a.distAnt2(a.pursue)
			if dist2m > dist2Max*32 {
				a.pursue = nil
			} else if dist2m < dist2Max/4 {
				a.pursue.decrLife(ns, 1)
				index := a.getDirection(a.pursue.X, a.pursue.Y)
				a.entries[visionNb*3+index] = ((dist2Max - dist2m) / dist2Max)
				return true
			} else {
				index := a.getDirection(a.pursue.X, a.pursue.Y)
				a.entries[visionNb*3+index] = ((dist2Max - dist2m) / dist2Max)
				return true
			}
		} else {
			a.pursue = nil
		}
	}
	a.timeWithoutHostile++
	var antMin *Ant
	dist2m := dist2Max
	if a.AntType == 1 {
		a.speed = a.regularSpeed
		dist2m = dist2m * 8
	}
	for _, nest := range ns.nests {
		if nest.id != a.nest.id {
			for _, ant := range nest.ants {
				if ant.Life > 0 {
					dist2 := a.distAnt2(ant)
					if dist2 < dist2m {
						antMin = ant
						dist2m = dist2
					}
				}
			}
		}
	}
	if antMin != nil {
		a.soldierInitCounter = 0
		a.timeWithoutHostile = 0
		if a.AntType == 1 {
			a.pursue = antMin
			a.speed = a.maxSpeed
		}
		if dist2m < dist2Max/16 {
			a.Fight = true
			antMin.Fight = true
			if a.AntType == 1 {
				antMin.decrLife(ns, 2)
				a.decrLife(ns, 1)
			}
		}
		if a.AntType == 0 {
			if rand.Float64() < 0.6 && (a.X-a.nest.x)*(a.X-a.nest.x)+(a.Y-a.nest.y)*(a.Y-a.nest.y) > 4000 {
				a.printf(ns, "current ant panic mode\n")
				a.panic = true
				a.Fight = false
				if a.food != nil {
					a.dropFood(ns)
				}
			}
		}
		index := a.getDirection(antMin.X, antMin.Y)
		a.entries[visionNb*3+index] = ((dist2Max - dist2m) / dist2Max)
		return true
	}
	return false
}

func (a *Ant) updateEntriesForFriendAnts(ns *Nests) bool {
	dist2Max := a.vision * a.vision
	dist2m := dist2Max
	var antMin *Ant
	if a.AntType == 0 {
		for _, ant := range a.nest.ants {
			if ant.Life > 0 && ant.AntType == 0 && ant != a && ant.food == nil {
				dist2 := a.distAnt2(ant)
				if dist2 < dist2Max/16 {
					a.statContact.incr()
					a.Contact = true
				}
				if dist2 < dist2m {
					antMin = ant
					dist2m = dist2
				}
			}
		}
	}
	if a.AntType == 1 && a.timeWithoutHostile > 10000 {
		for _, ant := range a.nest.ants {
			if ant.Life > 0 && ant != a && ant.food == nil {
				dist2 := a.distAnt2(ant)
				if dist2 < dist2Max/16 {
					a.statContact.incr()
					a.Contact = true
				}
				if dist2 < dist2m {
					antMin = ant
					dist2m = dist2
				}
			}
		}
	}

	if antMin != nil {
		//a.entryMode = 1
		index := a.getDirection(antMin.X, antMin.Y)
		a.entries[index] = ((dist2Max - dist2m) / dist2Max)
		return true
	}
	return false
}

func (a *Ant) getDirection(x float64, y float64) int {
	ang := math.Atan2(x-a.X, y-a.Y)
	if ang < 0 {
		ang = 2*math.Pi + ang
	}
	index := int(ang*float64(visionNb)/2.0/math.Pi + 0.000001)
	if index >= visionNb {
		index = index - visionNb
	}
	return index
}

func (a *Ant) computeHappiness(ns *Nests) {
	a.lastHappiness = a.happiness
	a.happiness = 0
	if a.AntType == 1 {
		//hostile ants
		for ii := visionNb * 3; ii < visionNb*4; ii++ {
			a.happiness += a.entries[ii]
		}
		//friend ants
		if a.happiness == 0 {
			for ii := 0; ii < visionNb; ii++ {
				a.happiness -= a.entries[ii]
			}
		}
		return
	}
	//foods
	for ii := visionNb; ii < visionNb*2; ii++ {
		a.happiness += a.entries[ii]
	}
	//Food Pheromones
	if a.happiness == 0 {
		for ii := visionNb * 2; ii < visionNb*3; ii++ {
			a.happiness += a.entries[ii]
		}
	}
	//friend ants
	if a.happiness == 0 {
		for ii := 0; ii < visionNb; ii++ {
			a.happiness -= a.entries[ii]
		}
	}
}

func (a *Ant) getDirIndex(nn int) int {
	if nn >= outNb {
		return nn - outNb
	}
	if nn < 0 {
		return nn + outNb
	}
	return nn
}

func (a *Ant) moveOnOut(ns *Nests) {
	//for now the nest return is hard coded
	if a.food != nil || a.panic {
		a.moveToNest(ns)
		return
	}
	if a.AntType == 1 && a.soldierInitCounter > 0 {
		a.soldierInitCounter--
		a.X += a.dx * a.speed
		a.Y += a.dy * a.speed
	} else {
		angle := (math.Pi * 2 * float64(a.Direction)) / float64(outNb) //+ math.Pi/2
		a.X += math.Sin(angle) * a.speed
		a.Y += math.Cos(angle) * a.speed
	}

	if a.X < ns.xmin {
		//a.x = ns.xmax
		a.X = ns.xmin
		a.Direction = 1 + int(rand.Intn(outNb/4+1))
	} else if a.Y < ns.ymin {
		//a.y = ns.ymax
		a.Y = ns.ymin
		a.Direction = 7 + int(rand.Intn(outNb/4+1))
	} else if a.X > ns.xmax {
		//a.x = ns.xmin
		a.X = ns.xmax
		a.Direction = 5 + int(rand.Intn(outNb/4+1))
	} else if a.Y > ns.ymax {
		//a.y = ns.ymin
		a.Y = ns.ymax
		a.Direction = 3 + int(rand.Intn(outNb/4+1))
	}
	if a.Direction >= outNb {
		a.Direction = a.Direction - outNb
	}
}

func (a *Ant) moveToNest(ns *Nests) {
	speed := a.speed
	if a.panic {
		speed = speed * 2
	}
	dd := math.Sqrt(float64((a.nest.x-a.X)*(a.nest.x-a.X) + (a.nest.y-a.Y)*(a.nest.y-a.Y)))
	dx := (a.nest.x - a.X) / dd
	dy := (a.nest.y - a.Y) / dd
	a.Direction = a.getDirection(a.nest.x, a.nest.y)
	a.X += dx * speed
	a.Y += dy * speed
	if a.food != nil {
		if a.nest.id == 1 {
			a.food.X = a.X
			a.food.Y = a.Y
		} else {
			a.food.X = a.X + 1
			a.food.Y = a.Y + 1
		}
		a.pheromoneDelay--
		if a.pheromoneDelay <= 0 {
			a.printf(ns, "add food pheromone\n")
			a.pheromoneCount++
			a.nest.addPheromone(a.X, a.Y, a.pheromoneCount)
			a.pheromoneDelay = a.nest.parameters.pheromoneAntDelay
		}
	}
	if (a.nest.x-a.X)*(a.nest.x-a.X)+(a.nest.y-a.Y)*(a.nest.y-a.Y) < 4000 {
		direc := a.Direction + outNb/2
		if direc >= outNb {
			direc = direc - outNb
		}
		if a.panic {
			a.nest.addSoldier(ns, a.X, a.Y, -dx, -dy, direc)
			a.panic = false
		}
		if a.food != nil {
			a.nest.ressource += 4
			if len(ns.foodGroups) > 0 {
				if ns.foodRenew {
					a.food.renew()
				}
			}
			a.Direction = direc
			a.nest.addWorker(ns, a.X, a.Y, direc)
			a.dropFood(ns)
		}
	}
}

func (a *Ant) train(ns *Nests) bool {
	if a.lastDecision < 0 || a.lastEntryMode != a.entryMode {
		return false
	}
	//fmt.Printf("%d  entries: %v\n", a.id, a.lastEntries)
	ins, ok := a.preparedEntries(a.lastEntries)
	if !ok {
		return false
	}
	outs := a.network.Propagate(ins, true)
	a.setOuts(a.lastDecision)
	if ns.log && a.ID == ns.selected {
		trainResult := a.computeTrainResult(a.outs, outs)
		a.printf(ns, "train %s => %v\n", a.displayList(ns, ins, "%.0f"), a.outs)
		a.printf(ns, "outs: %s result=%f\n", a.displayList(ns, outs, "%.3f"), trainResult)
	}
	if a.ID == ns.selected {
		ns.addSample(ins, a.outs)
	}
	a.network.BackPropagate(a.outs)
	return true
}

func (a *Ant) fadeLastDecision(ns *Nests) bool {
	if a.lastDecision == -1 || a.entryMode != a.lastEntryMode {
		return false
	}
	ins, ok := a.preparedEntries(a.lastEntries)
	if !ok {
		return false
	}
	outs := a.network.Propagate(ins, true)
	a.setOutsFaded(a.lastDecision)
	if ns.log && a.ID == ns.selected {
		trainResult := a.computeTrainResult(a.outs, outs)
		a.printf(ns, "fade %s => %v\n", a.displayList(ns, ins, "%.0f"), a.outs)
		a.printf(ns, "outs: %s result=%.5f\n", a.displayList(ns, outs, "%.3f"), trainResult)
		ns.addSample(ins, a.outs)
	}
	a.network.BackPropagate(a.outs)
	a.statFade.incr()
	return true
}

func (a *Ant) preparedEntries(list []float64) ([]float64, bool) {
	ret := make([]float64, len(list), len(list))
	isNull := true
	for ii, val := range list {
		ret[ii] = val
		if val > 0 {
			ret[ii] = 1 //val
			isNull = false
		}
	}
	return ret, !isNull
}

func (a *Ant) setOuts(direction int) {
	for ii := range a.outs {
		a.outs[ii] = 0
	}
	a.outs[direction] = 1
}

func (a *Ant) setOutsFaded(lastDecision int) {
	for ii := range a.outs {
		a.outs[ii] = 0.2
	}
	a.outs[lastDecision] = 0
}

func (a *Ant) dropFood(ns *Nests) {
	if a.food != nil {
		if ns.foodRenew {
			a.food.carried = false
		}
		a.food = nil
	}
}

func (a *Ant) carryFood(f *Food) {
	if !f.carried {
		a.food = f
		f.carried = true
	}
}
