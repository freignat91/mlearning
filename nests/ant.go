package nests

import (
	"math"
	"math/rand"

	"github.com/freignat91/mlearning/network"
)

var visionNb = 8
var inNb = 8 * 4 //8 for ants, 8 for foods, 8 for pheromones, 8 for hostile
var outNb = 8
var modes = []string{"free", "spread", "found food", "piste up food", "hostile"}
var deltaBase = 1 // number of train to do for the first partition 0 or 1
var deltaCoef = 1
var updateNetwork = true

// Ant .
type Ant struct {
	ID        int     `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Direction int     `json:"direction"`
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
	lastDecision       int
	lastLastDecision   int
	lastResult         bool
	lastDecisionRandom bool
	panic              bool
	statDecision       *Stats
	statGoodDecision   *Stats
	statReinforce      *Stats
	statFade           *Stats
	gRate              float64
	dirMap             []int
	dirCount           int
	food               *Food
	pheromoneDelay     int
	pheromoneCount     int
	lastPheromone      int
	lastPheromoneCount int
	lastEntryMode      int
	timeWithoutHostile int
	trained            bool
	pursue             *Ant
	mode               int
	lastMode           int
	//
	decTmp []int
	//happinnessDeltaMax    float64
	//happinnessDeltaMaxTmp float64
	happinessDeltaSum     [5]float64
	happinessDeltaNumber  [5]int
	averageHappinessDelta [5]float64
	//debug
	index   int
	tmpFood *Food
	lastx   float64
	lasty   float64
}

func newAnt(ns *Nests, n *Nest, antType int) *Ant {
	n.antIDCounter++
	ant := &Ant{
		ID:               n.antIDCounter,
		X:                n.x + 20.0 - rand.Float64()*40,
		Y:                n.y + 20.0 - rand.Float64()*40,
		AntType:          antType,
		nest:             n,
		vision:           30,
		Direction:        int(rand.Int31n(int32(outNb))),
		entries:          make([]float64, inNb, inNb),
		lastEntries:      make([]float64, inNb, inNb),
		outs:             make([]float64, outNb, outNb),
		lastDecision:     -1,
		statDecision:     newStats(n.statDecision),
		statGoodDecision: newStats(n.statGoodDecision),
		statReinforce:    newStats(n.statReinforce),
		statFade:         newStats(n.statFade),
		dirMap:           make([]int, outNb, outNb),
		decTmp:           make([]int, outNb/2, outNb/2),
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
	if updateNetwork {
		if a.AntType == 0 {
			if a.nest.bestWorker != nil && rand.Float64() < a.nest.parameters.chanceToGetTheBestNetworkCopy {
				net, err := a.nest.bestWorker.network.Copy()
				if err == nil {
					a.network = net
					return
				}
			}
		} else {
			if a.nest.bestSoldier != nil && rand.Float64() < a.nest.parameters.chanceToGetTheBestNetworkCopy {
				net, err := a.nest.bestSoldier.network.Copy()
				if err == nil {
					a.network = net
					return
				}
			}
		}
	}
	var defnet []int
	if rand.Float64() < 0.5 || true {
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

func (a *Ant) nextTime(ns *Nests, update bool) {
	if !a.tickLife(ns) {
		return
	}
	a.lastLastDecision = a.lastDecision
	a.lastDecision = a.Direction
	a.displayInfo1(ns)
	a.moveOnOut(ns)
	if a.panic {
		return
	}
	a.updateEntries(ns)
	a.computeHappiness(ns)
	a.displayInfo2(ns)
	if ns.log && a == ns.selected {
		a.printf(ns, "mode: %d entries: %s\n", a.mode, a.displayList(ns, a.entries, "%.2f"))
		a.printf(ns, "happiness=%.3f\n", a.happiness)
	}
	if a.happiness == a.lastHappiness && a.happiness >= 0 {
		a.printf(ns, "decision: no need\n")
	} else if a.happiness < a.lastHappiness {
		a.printf(ns, "bad last decision: %d\n", a.lastDecision)
		if rand.Float64() > 1-(a.gRate-5)/100 && (a.lastResult || a.lastDecision != a.lastLastDecision) {
			a.decide(ns)
			a.printf(ns, "decision using network: %d\n", a.Direction)
			a.lastDecisionRandom = false
			a.statDecision.incr()
		} else {
			a.Direction = int(rand.Int31n(int32(outNb)))
			//a.Direction = a.index + 4
			//if a.Direction >= outNb {
			//a.Direction = a.Direction - outNb
			//}
			a.dirMap[a.Direction]++
			a.lastDecisionRandom = true
			a.statDecision.incr()
			a.printf(ns, "decision random: %d\n", a.Direction)
		}
		a.lastResult = false
	} else if a.lastDecision != -1 && a.mode == a.lastMode {
		a.statGoodDecision.incr()
		a.lastResult = true
		a.printf(ns, "decision: no need\n")
		delta := (a.happiness - a.lastHappiness)
		if delta > 0 && (delta < a.averageHappinessDelta[a.mode]*3 || a.averageHappinessDelta[a.mode] == 0) {
			a.happinessDeltaSum[a.mode] += delta
			a.happinessDeltaNumber[a.mode]++
		}
		a.printf(ns, "positive decision %d, delta=%.4f average=%.4f\n", a.lastDecision, delta, a.averageHappinessDelta)
		_, ok := a.train(ns)
		if ok {
			a.printf(ns, "positive decision reinforced\n")
			a.statReinforce.incr()
			a.lastDecision = -1
		}
	}
	if update {
		a.update(ns)
	}
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

func (a *Ant) update(ns *Nests) {
	if updateNetwork && ns.timeRef%5000 == 0 {
		bestAnt := a.nest.bestWorker
		if a.AntType == 1 {
			bestAnt = a.nest.bestSoldier
		}
		if a.dirCount < bestAnt.dirCount-3 || (a.dirCount <= bestAnt.dirCount && a.gRate < bestAnt.gRate-10) {
			net, err := bestAnt.network.Copy()
			if err == nil {
				a.network = net
				a.lastDecision = -1
				a.printf(ns, "update network with the best one: %v\n", a.network.Getdef())
			}
		}
	}
}

func (a *Ant) decide(ns *Nests) bool {
	if a.food != nil {
		return false
	}
	ins, ok := a.preparedEntries(a.entries)
	if !ok {
		a.printf(ns, "bad entries: %s\n", a.displayList(ns, ins, "%.3f"))
		return false
	}
	a.outs = a.network.Propagate(a.entries, true)
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
	a.Direction = direction
	a.dirMap[direction]++
	return true
}

func (a *Ant) commitPeriodStats(ns *Nests) {
	if ns.stopped {
		return
	}
	nb := 0
	for _, val := range a.dirMap {
		if val > 0 {
			nb++
		}
	}
	a.dirCount = nb
	a.dirMap = make([]int, outNb, outNb)
	//
	for ii := 1; ii < 5; ii++ {
		a.averageHappinessDelta[ii] = 0
		if a.happinessDeltaNumber[ii] > 0 {
			a.averageHappinessDelta[ii] = a.happinessDeltaSum[ii] / float64(a.happinessDeltaNumber[ii])
		}
		a.happinessDeltaSum[ii] = 0
		a.happinessDeltaNumber[ii] = 0
	}
	if a.statDecision.value != 0 {
		a.statDecision.push()
		a.statGoodDecision.push()
		a.statReinforce.push()
		a.statFade.push()
		a.gRate = float64(a.statGoodDecision.cumul) * 100.0 / float64(a.statDecision.cumul)
		if a.gRate > 100 {
			a.gRate = 100
		}
	}
}

func (a *Ant) updateEntries(ns *Nests) {
	if ns.timeRef%100 == 0 {
		a.Fight = false
	}
	for ii := range a.entries {
		a.lastEntries[ii] = a.entries[ii]
		a.entries[ii] = 0
	}
	if a.food != nil {
		a.mode = 0
		a.printf(ns, "Carry food no entries\n")
		return
	}
	a.lastMode = a.mode
	a.mode = 0
	if a.updateEntriesForHostileAnts(ns) {
		a.mode = 4
		return
	}
	if a.AntType == 0 {
		if a.updateEntriesForFoods(ns) {
			a.mode = 2
			return
		}
		if a.updateEntriesForPheromones(ns) {
			a.mode = 3
			return
		}
	}
	a.updateEntriesForFriendAnts(ns)
	a.mode = 1
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
		if dist2m < 5 {
			a.carryFood(foodMin)
			a.pheromoneCount = 0
			return true
		}
		a.tmpFood = foodMin
		index := a.getDirection(ns, foodMin.X, foodMin.Y)
		a.entries[visionNb+index] = ((dist2Max - dist2m) / dist2Max)
		return true
	}
	return false
}

func (a *Ant) updateEntriesForPheromones(ns *Nests) bool {
	//minLevel := a.nest.parameters.pheromoneLevel + 1
	minLevel := 1000000
	dist2Max := a.vision * a.vision * 9

	dist2m := dist2Max
	var pheMin *Pheromone
	for _, phe := range a.nest.pheromones {
		if phe.Level > 0 {
			dist2 := (phe.X-a.X)*(phe.X-a.X) + (phe.Y-a.Y)*(phe.Y-a.Y)
			if dist2 < dist2Max && phe.id < minLevel {
				minLevel = phe.id
				index := a.getDirection(ns, phe.X, phe.Y)
				//if ns.selected == a {
				//	fmt.Printf("Pheromone direction: %d\n", index)
				//}
				pheMin = phe
				a.entries[visionNb*2+index] += (a.nest.parameters.pheromoneLevel - float64(phe.Level)) * ((dist2Max - dist2) / dist2Max)
				if dist2 < dist2m {
					pheMin = phe
					dist2m = dist2
				}
			}
		}
	}
	a.printf(ns, "pheromone: %+v\n", pheMin)
	if pheMin != nil {

		if a.lastPheromone == pheMin.id {
			a.lastPheromoneCount++
		} else {
			a.lastPheromoneCount = 0
		}
		a.lastPheromone = pheMin.id
		if a.lastPheromoneCount > 200 {
			a.printf(ns, "same pheromone too much time: ignored\n")
			for ii := visionNb * 2; ii < visionNb*3; ii++ {
				a.entries[ii] = 0
			}
		}
		return true
	}
	return false
}

func (a *Ant) updateEntriesForHostileAnts(ns *Nests) bool {
	if a.AntType == 0 && !panicMode {
		return false
	}
	dist2Max := a.vision * a.vision
	dist2Contact := dist2Max / 8
	if a.AntType == 1 {
		dist2Max = dist2Max * 16
	}
	//pursue mode
	if a.AntType == 1 && a.pursue != nil {
		if a.pursue.Life > 0 {
			a.printf(ns, "Pursue ant: %d\n", a.pursue.ID)
			dist2m := a.distAnt2(a.pursue)
			if dist2m > dist2Max*4 {
				a.pursue = nil
			} else if dist2m < dist2Contact {
				a.pursue.decrLife(ns, 2)
				index := a.getDirection(ns, a.pursue.X, a.pursue.Y)
				a.entries[visionNb*3+index] = ((dist2Max - dist2m) / dist2Max)
				return true
			} else {
				index := a.getDirection(ns, a.pursue.X, a.pursue.Y)
				a.entries[visionNb*3+index] = ((dist2Max*4 - dist2m) / dist2Max / 4)
				return true
			}
		} else {
			a.pursue = nil
		}
	}
	a.speed = a.regularSpeed
	a.timeWithoutHostile++
	var antMin *Ant
	dist2m := dist2Max
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
	a.printf(ns, "closest hostile: %+v\n", antMin)
	if antMin != nil {
		a.soldierInitCounter = 0
		a.timeWithoutHostile = 0
		if a.AntType == 1 {
			a.pursue = antMin
			a.speed = a.maxSpeed
		}
		if dist2m < dist2Contact {
			a.Fight = true
			antMin.Fight = true
			if a.AntType == 1 {
				antMin.decrLife(ns, 2)
				a.decrLife(ns, 1)
			}
		}
		if a.AntType == 0 {
			if rand.Float64() < 0.01 && (a.X-a.nest.x)*(a.X-a.nest.x)+(a.Y-a.nest.y)*(a.Y-a.nest.y) > 4000 {
				a.printf(ns, "current ant panic mode\n")
				a.panic = true
				a.Fight = false
				if a.food != nil {
					a.dropFood(ns)
				}
			}
		}
		index := a.getDirection(ns, antMin.X, antMin.Y)
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
				}
				if dist2 < dist2m {
					antMin = ant
					dist2m = dist2
				}
			}
		}
	}
	if a.AntType == 1 && a.timeWithoutHostile > 5000 {

		for _, ant := range a.nest.ants {
			if ant.Life > 0 && ant != a && ant.food == nil {
				dist2 := a.distAnt2(ant)
				if dist2 < dist2m {
					antMin = ant
					dist2m = dist2
				}
			}
		}

		for _, ant := range a.nest.ants {
			if ant.Life > 0 && ant != a && ant.food == nil {
				dist2 := a.distAnt2(ant)
				if dist2 < dist2Max {
					antMin = ant
					dist2m = dist2
					index := a.getDirection(ns, antMin.X, antMin.Y)
					a.entries[index] += ((dist2Max - dist2m) / dist2Max)
				}
			}
		}
	}
	a.printf(ns, "closest friend: %+v\n", antMin)
	if antMin != nil {
		index := a.getDirection(ns, antMin.X, antMin.Y)
		a.index = index
		a.entries[index] = ((dist2Max - dist2m) / dist2Max)
		return true
	}
	return false
}

func (a *Ant) getDirection(ns *Nests, x float64, y float64) int {
	ang := math.Atan2(x-a.X, y-a.Y)
	if ang < 0 {
		ang = 2*math.Pi + ang
	}
	index := int(ang/(math.Pi*2.0/float64(outNb))) + 1
	if index >= outNb {
		return index - outNb
	}
	return index
}

func (a *Ant) computeHappiness(ns *Nests) {
	a.lastHappiness = a.happiness
	a.happiness = 0
	for ii, val := range a.entries {
		if ii >= 0 && ii < visionNb {
			a.happiness -= val
		} else {
			a.happiness += val
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
	a.lastx = a.X
	a.lasty = a.Y
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

	max := 2.0
	if a.X < ns.xmin*max {
		//a.X = ns.xmax
		a.X = ns.xmin * max
		a.Direction = 1 + int(rand.Intn(outNb/4+1))
	} else if a.Y < ns.ymin*max {
		//a.Y = ns.ymax
		a.Y = ns.ymin * max
		a.Direction = 7 + int(rand.Intn(outNb/4+1))
	} else if a.X > ns.xmax*max {
		//a.X = ns.xmin
		a.X = ns.xmax * max
		a.Direction = 5 + int(rand.Intn(outNb/4+1))
	} else if a.Y > ns.ymax*max {
		//a.Y = ns.ymin
		a.Y = ns.ymax * max
		a.Direction = 3 + int(rand.Intn(outNb/4+1))
	}
	if a.Direction >= outNb {
		a.Direction = a.Direction - outNb
	}
}

func (a *Ant) moveToNest(ns *Nests) {
	speed := a.speed
	if a.panic {
		a.Fight = false
		speed = speed * 2
	}
	dd := math.Sqrt(float64((a.nest.x-a.X)*(a.nest.x-a.X) + (a.nest.y-a.Y)*(a.nest.y-a.Y)))
	dx := (a.nest.x - a.X) / dd
	dy := (a.nest.y - a.Y) / dd
	a.Direction = a.getDirection(ns, a.nest.x, a.nest.y)
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
	if (a.nest.x-a.X)*(a.nest.x-a.X)+(a.nest.y-a.Y)*(a.nest.y-a.Y) < 1600 {
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
				if foodRenew {
					a.food.renew()
				}
			}
			a.Direction = direc
			a.nest.addWorker(ns, a.X, a.Y, direc)
			a.dropFood(ns)
		}
	}
}

func (a *Ant) train(ns *Nests) (int, bool) {
	if a.lastDecision < 0 {
		return 0, false
	}
	if a.trained {
		return 0, true
	}
	nb := a.getNbTrain(ns)
	if nb == 0 {
		return 0, false
	}

	ins, ok := a.preparedEntries(a.lastEntries)
	if !ok {
		return 0, false
	}
	if !a.lastDecisionRandom {
		return 0, true
	}
	//train a much as the decision appears good concidering delta happiness stats
	for ii := 0; ii < nb; ii++ {
		a.network.Propagate(ins, false)
		a.setOuts(a.lastDecision)
		if a == ns.selected {
			ns.addSample(ins, a.outs)
		}
		a.network.BackPropagate(a.outs)
	}
	return nb, true
}

func (a *Ant) getNbTrain(ns *Nests) int {
	delta := a.happiness - a.lastHappiness
	var ret int
	if a.averageHappinessDelta[a.mode] <= 0 {
		ret = 0
	} else if delta < a.averageHappinessDelta[a.mode]/2 {
		ret = deltaBase
	} else if delta < a.averageHappinessDelta[a.mode] {
		ret = deltaBase
	} else if delta < a.averageHappinessDelta[a.mode]*1.5 {
		ret = deltaBase + 1
	} else {
		ret = deltaBase + 1
	}
	//a.printf(ns, "Positive decision, delta=%.3f average=%.3f max=%.3f nbTrain=%d\n", delta, a.averageHappinessDelta, a.happinnessDeltaMax, ret)
	a.printf(ns, "Positive decision, delta=%.3f average=%.3f nbTrain=%d\n", delta, a.averageHappinessDelta, ret)
	return ret * deltaCoef
}

func (a *Ant) fadeLastDecision(ns *Nests) bool {
	if a.lastDecision == -1 || a.trained {
		return false
	}
	ins, ok := a.preparedEntries(a.lastEntries)
	if !ok {
		return false
	}
	outs := a.network.Propagate(ins, true)
	a.setOutsFaded(a.lastDecision)
	if ns.log && a == ns.selected {
		a.computeTrainResult(a.outs, outs)
		a.printf(ns, "fade last decision: %d\n", a.lastDecision)
		ns.addSample(ins, a.outs)
	}
	a.network.BackPropagate(a.outs)
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
		if foodRenew {
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

func (a *Ant) displayInfo1(ns *Nests) {
	if ns.log && ns.selected == a {
		a.printf(ns, "-----------------------------------------------------\n")
		ggRate := float64(a.statReinforce.scumul) * 100.0 / float64(a.statDecision.scumul)
		a.printf(ns, "%d:[%d] type=%d reinforce:%d decision:%d period:good=%.2f%%) global:good=%.2f%%)\n", ns.timeRef, a.ID, a.AntType, a.statReinforce.cumul, a.statDecision.cumul, a.gRate, ggRate)
		a.printf(ns, "network=%v hapiness=%.5f direction: %d last decision: %d (%d) result: %t\n", a.network.Getdef(), a.happiness, a.Direction, a.lastDecision, a.lastLastDecision, a.lastResult)
	}
}

func (a *Ant) displayInfo2(ns *Nests) {
	if ns.log && ns.selected == a {
		a.printf(ns, "delta=%.4f average=%.4f\n", a.happiness-a.lastHappiness, a.averageHappinessDelta)
	}
}

func (a *Ant) getModeToString() string {
	if a.panic {
		return "panic"
	}
	if a.food != nil {
		return "carry food"
	}
	if a.AntType == 0 {
		if a.Fight {
			return "attacked"
		}
	} else {
		if a.Fight {
			return "attack"
		}
		if a.pursue != nil {
			return "pursue"
		}
	}
	return modes[a.mode]
}
