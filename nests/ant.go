package nests

import (
	"math"
	"math/rand"

	"github.com/freignat91/mlearning/network"
)

var visionNb = 8
var inNb = 8 * 3 //8 for ants, 8 for foods, 8 for pheromones
var outNb = 8

// Ant .
type Ant struct {
	id                 int
	nest               *Nest
	happiness          float64
	lastHappiness      float64
	x                  float64
	y                  float64
	speed              float64
	vision             float64
	maxSpeed           float64
	networkDef         []int
	network            *network.MLNetwork
	entries            []float64
	lastEntries        []float64
	outs               []float64
	lastDecision       int
	statTrain          *Stats
	statDecision       *Stats
	statReinforce      *Stats
	statFade           *Stats
	statNetwork        *Stats
	statContact        *Stats
	direction          int
	iner               int
	gRate              float64
	life               int64
	dead               bool
	contact            bool
	edge               bool
	dirMap             map[int]int
	dirCount           int
	distinctOuts       int
	carryFood          *Food
	pheromoneDelay     int
	pheromoneCount     int
	lastPheromone      int
	lastPheromoneCount int
	entryMode          int
	lastEntryMode      int
	//
	decTmp []int
}

func newAnt(ns *Nests, n *Nest, id int) (*Ant, error) {
	ant := &Ant{
		id:            id,
		nest:          n,
		maxSpeed:      0.3, //+ rand.Float64()/10,
		vision:        30,
		x:             n.x + 20.0 - rand.Float64()*40,
		y:             n.y + 20.0 - rand.Float64()*40,
		direction:     int(rand.Int31n(int32(outNb))),
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
	ant.speed = ant.maxSpeed
	ant.setNetwork(ns)
	if id == 1 {
		/*
			ant.x = ns.xmax / 2
			ant.y = ns.ymax / 2
			ant.direction = 0
			//ant.speed = 0
		*/
	} else if ant.id <= inNb+1 {
		/*
			ant.x = (ns.xmax / 2) + 25*math.Sin(math.Pi*2.0/float64(visionSize)*float64(ant.id-2))
			ant.y = (ns.ymax / 2) + 25*math.Cos(math.Pi*2.0/float64(visionSize)*float64(ant.id-2))
			ant.direction = ant.id - 2
		*/
	}
	//fmt.Printf("ant: %+v\n", ant)
	return ant, nil
}

func (a *Ant) setNetwork(ns *Nests) {
	var defnet []int
	if rand.Float64() < 0.5 {
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
	a.networkDef = defnet
	net, _ := network.NewNetwork(a.networkDef)
	a.network = net
}

func (a *Ant) getData() *AntData {
	data := AntData{
		ID:        a.id,
		X:         a.x,
		Y:         a.y,
		Direction: a.direction,
		Contact:   a.contact,
	}
	return &data
}

func (a *Ant) nextTime(ns *Nests) {
	a.displayInfo(ns)
	a.moveOnOut(ns)
	a.updateEntries(ns)
	a.computeHappiness(ns)
	a.printf(ns, "happiness=%.3f entries: %s\n", a.happiness, a.displayList(ns, a.entries, "%.2f"))
	if ns.random {
		a.direction = int(rand.Int31n(int32(outNb)))
	} else {
		if ns.log && ns.selected == a.id {
			a.printf(ns, "lastDecision:%d entrie: %s hapiness=%.5f delta=%.5f\n", a.lastDecision, a.displayList(ns, a.entries, "%.3f"), a.happiness, a.happiness-a.lastHappiness)
		}
		if a.iner < 0 {
			if a.happiness == a.lastHappiness && a.happiness >= 0 {
				a.iner = -1
				a.printf(ns, "decision: no need\n")
			} else if a.happiness <= a.lastHappiness {
				a.fadeLastDecision(ns)
				if a.decide(ns) {
					if ns.log && ns.selected == a.id {
						a.printf(ns, "decision using network: %d outs: %s\n", a.direction, a.displayList(ns, a.outs, "%.3f"))
					}
					a.statDecision.incr()
				} else {
					a.direction = int(rand.Int31n(int32(outNb)))
					a.lastDecision = -1
					a.printf(ns, "decision random: %d\n", a.direction)
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
		}
	}
	a.iner--
	a.update(ns)
}

func (a *Ant) decide(ns *Nests) bool {
	if ns.random {
		return false
	}
	if rand.Float64() < 0.1 {
		//return false
	}
	if a.carryFood != nil {
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
	a.direction = direction
	a.lastDecision = direction
	a.dirMap[direction]++
	return true
}

func (a *Ant) update(ns *Nests) {
	a.life++
	if a.life > 100000 && a.dirCount > 0 {
		a.life = 0
		if false { //a.dirCount < 3 {
			a.setNetwork(ns)
			a.statNetwork.incr()
			a.printf(ns, "recreate new random network: %v\n", a.networkDef)
		} else if a.dirCount < ns.bestAnt.dirCount-4 || (a.dirCount <= ns.bestAnt.dirCount && a.gRate < ns.bestAnt.gRate-10) {
			net, err := ns.bestAnt.network.Copy()
			if err == nil {
				a.network = net
				a.networkDef = net.Getdef()
				a.lastDecision = -1
				a.statNetwork.incr()
				a.printf(ns, "update network with the best one: %v\n", a.networkDef)
			}
		}
	}
}

func (a *Ant) displayInfo(ns *Nests) {
	if ns.timeRef%10000 == 0 {
		a.dirMap = make(map[int]int)
		a.dirCount = a.network.ComputeDistinctOut()
	}
	if ns.log && a.id == ns.selected {
		ggRate := float64(a.statReinforce.scumul) * 100.0 / float64(a.statDecision.scumul)
		a.printf(ns, "[%d] totTrain: %d train:%d reinforce:%d fade:%d decision:%d period:good=%.2f%%) global:good=%.2f%%)\n", a.id, a.statTrain.scumul, a.statTrain.cumul, a.statReinforce.cumul, a.statFade.cumul, a.statDecision.cumul, a.gRate, ggRate)
		a.printf(ns, "network=%v hapiness=%.5f move: %d\n", a.networkDef, a.happiness, a.direction)
	}
}

func (a *Ant) commitPeriodStats(ns *Nests) {
	if ns.stopped {
		return
	}
	a.statTrain.push()
	a.statDecision.push()
	a.statReinforce.push()
	a.statFade.push()
	a.statNetwork.push()
	a.statContact.push()
	if a.statDecision.cumul == 0 {
		a.gRate = 0
	} else {
		a.gRate = float64(a.statReinforce.cumul) * 100.0 / float64(a.statDecision.cumul)
	}
}

func (a *Ant) updateEntries(ns *Nests) {
	for ii := range a.entries {
		a.lastEntries[ii] = a.entries[ii]
		a.entries[ii] = 0
	}
	if a.carryFood != nil {
		return
	}
	if a.updateEntriesForFoods(ns) {
		return
	}
	if a.updateEntriesForPheromones(ns) {
		return
	}
	a.updateEntriesForFriendAnts(ns)

}

func (a *Ant) updateEntriesForFoods(ns *Nests) bool {
	dist2Max := a.vision * a.vision
	a.contact = false
	a.lastEntryMode = a.entryMode
	//search food
	dist2m := dist2Max
	var foodMin *Food
	for _, food := range ns.foods {
		if !food.carried {
			dist2 := a.distFood2(food)
			if dist2 < dist2m {
				foodMin = food
				dist2m = dist2
			}
		}
	}
	a.printf(ns, "closest food: %+v\n", foodMin)
	if foodMin != nil {
		a.entryMode = 2
		if dist2m < 4 {
			a.carryFood = foodMin
			foodMin.carried = true
			a.pheromoneCount = 0
			return true
		}
		ang := math.Atan2(foodMin.X-a.x, foodMin.Y-a.y)
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
			dist2 := a.distPhe2(phe)
			if dist2 < dist2Max*1.5*1.5 && phe.id < minLevel {
				dist2m = dist2
				pheMin = phe
				minLevel = phe.id
			}
		}
	}
	if pheMin != nil {
		a.entryMode = 3
		if a.lastPheromone == pheMin.id {
			a.lastPheromoneCount++
		} else {
			a.lastPheromoneCount = 0
		}
		a.lastPheromone = pheMin.id
		if a.lastPheromoneCount < 100 {
			ang := math.Atan2(pheMin.X-a.x, pheMin.Y-a.y)
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

func (a *Ant) updateEntriesForFriendAnts(ns *Nests) bool {
	dist2Max := a.vision * a.vision
	dist2m := dist2Max
	var antMin *Ant
	for _, ant := range a.nest.ants {
		if ant != a && ant.carryFood == nil {
			dist2 := a.distAnt2(ant)
			if dist2 < dist2Max/4 {
				a.statContact.incr()
				a.contact = true
			}
			if dist2 < dist2m {
				antMin = ant
				dist2m = dist2
			}
		}
	}
	if antMin != nil {
		a.entryMode = 1
		ang := math.Atan2(antMin.x-a.x, antMin.y-a.y)
		if ang < 0 {
			ang = 2*math.Pi + ang
		}
		index := int(ang*float64(visionNb)/2.0/math.Pi + 0.000001)
		if index >= visionNb {
			index = index - visionNb
		}
		//a.printf(ns, "find %d angle=%0.2f degres=%0.2f index=%d\n", ant.id, ang, ang*180/math.Pi, index)
		a.entries[index] = ((dist2Max - dist2m) / dist2Max)
		return true
	}
	return false
}

func (a *Ant) computeHappiness(ns *Nests) {
	if a.edge {
		a.happiness = -.1
		return
	}
	a.lastHappiness = a.happiness
	a.happiness = 0
	//if a.carryFood != nil {
	//	a.happiness = float64((ns.xmax-ns.xmin)*(ns.xmax-ns.xmin)+(ns.ymax-ns.ymin)*(ns.ymax-ns.ymin)) / float64((a.nest.x-a.x)*(a.nest.x-a.x)+(a.nest.y-a.y)*(a.nest.y-a.y))
	//	return
	//}
	for ii := visionNb; ii < visionNb*2; ii++ {
		a.happiness += a.entries[ii]
	}
	if a.happiness == 0 {
		for ii := visionNb * 2; ii < visionNb*3; ii++ {
			a.happiness += a.entries[ii]
		}
	}
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
	if a.carryFood != nil {
		dd := math.Sqrt(float64((a.nest.x-a.x)*(a.nest.x-a.x) + (a.nest.y-a.y)*(a.nest.y-a.y)))
		a.x += (a.nest.x - a.x) / dd * a.speed
		a.y += (a.nest.y - a.y) / dd * a.speed
		a.carryFood.X = a.x
		a.carryFood.Y = a.y
		a.pheromoneDelay--
		if a.pheromoneDelay <= 0 {
			a.printf(ns, "add pheromone\n")
			a.pheromoneCount++
			a.nest.addPheromone(a.x, a.y, a.pheromoneCount)
			a.pheromoneDelay = a.nest.parameters.pheromoneAntDelay
		}
		if (a.nest.x-a.x)*(a.nest.x-a.x)+(a.nest.y-a.y)*(a.nest.y-a.y) < 4000 {
			a.nest.statFood.incr()
			if len(ns.foodGroups) > 0 {
				if ns.foodRenew {
					a.carryFood.renew()
					a.carryFood.carried = false
				}
			}
			a.carryFood = nil
		}
		return
	}
	angle := (math.Pi * 2 * float64(a.direction)) / float64(outNb) //+ math.Pi/2
	a.x += math.Sin(angle) * a.speed
	a.y += math.Cos(angle) * a.speed

	a.edge = false
	if a.x < ns.xmin {
		//a.x = ns.xmax
		a.x = ns.xmin
		a.direction = 1 + int(rand.Int31n(int32(outNb/3)))
		a.edge = true
	} else if a.y < ns.ymin {
		//a.y = ns.ymax
		a.y = ns.ymin
		a.direction = 10 + int(rand.Int31n(int32(outNb/3)))
		a.edge = true
	} else if a.x > ns.xmax {
		//a.x = ns.xmin
		a.x = ns.xmax
		a.direction = 7 + int(rand.Int31n(int32(outNb/3)))
		a.edge = true
	} else if a.y > ns.ymax {
		//a.y = ns.ymin
		a.y = ns.ymax
		a.direction = 4 + int(rand.Int31n(int32(outNb/3)))
		a.edge = true
	}
	if a.direction >= outNb {
		a.direction = a.direction - outNb
	}
}

func (a *Ant) train(ns *Nests) bool {
	if ns.random || a.lastDecision < 0 || a.lastEntryMode != a.entryMode {
		return false
	}
	//fmt.Printf("%d  entries: %v\n", a.id, a.lastEntries)
	ins, ok := a.preparedEntries(a.lastEntries)
	if !ok {
		return false
	}
	outs := a.network.Propagate(ins, true)
	a.setOuts(a.lastDecision)
	//a.setOuts(a.direction)
	//a.setOuts(a.soluce())
	if ns.log && a.id == ns.selected {
		trainResult := a.computeTrainResult(a.outs, outs)
		a.printf(ns, "train %s => %v\n", a.displayList(ns, ins, "%.0f"), a.outs)
		a.printf(ns, "outs: %s result=%f\n", a.displayList(ns, outs, "%.3f"), trainResult)
	}
	if a.id == ns.selected {
		ns.addSample(ins, a.outs)
	}
	a.network.BackPropagate(a.outs)
	return true
}

func (a *Ant) fadeLastDecision(ns *Nests) bool {
	if ns.random || a.lastDecision == -1 || a.entryMode != a.lastEntryMode {
		return false
	}
	ins, ok := a.preparedEntries(a.lastEntries)
	if !ok {
		return false
	}
	outs := a.network.Propagate(ins, true)
	a.setOutsFaded(a.lastDecision)
	if ns.log && a.id == ns.selected {
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

func (a *Ant) soluce() int {
	iim := 0
	val := 0.0
	for ii, in := range a.lastEntries {
		if in > val {
			val = in
			iim = ii
		}
	}
	iim += inNb / 2
	if iim >= inNb {
		iim = iim - inNb
	}
	return iim
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
