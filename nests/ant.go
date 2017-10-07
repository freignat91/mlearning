package nests

import (
	"math"
	"math/rand"

	"github.com/freignat91/mlearning/network"
)

var visionNb = 8
var inNb = 16 //8 for ants, 8 for foods
var outNb = 8

// Ant .
type Ant struct {
	id            int
	nestID        int
	hapiness      float64
	lastHapiness  float64
	x             float64
	y             float64
	speed         float64
	vision        float64
	maxSpeed      float64
	networkDef    []int
	network       *network.MLNetwork
	entries       []float64
	lastEntries   []float64
	outs          []float64
	lastDecision  int
	statTrain     *Stats
	statDecision  *Stats
	statReinforce *Stats
	statFade      *Stats
	statNetwork   *Stats
	statContact   *Stats
	direction     int
	iner          int
	gRate         float64
	life          int64
	dead          bool
	contact       bool
	edge          bool
	dirMap        map[int]int
	dirCount      int
	distinctOuts  int
	//
	decTmp []int
}

func newAnt(ns *Nests, n *Nest, id int) (*Ant, error) {
	ant := &Ant{
		id:            id,
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
	if ns.random {
		a.direction = int(rand.Int31n(int32(outNb)))
	} else {
		if ns.log && ns.selected == a.id {
			a.printf(ns, "lastDecision:%d entrie: %s hapiness=%.5f delta=%.5f\n", a.lastDecision, a.displayList(ns, a.entries), a.hapiness, a.hapiness-a.lastHapiness)
		}
		if a.iner < 0 {
			if a.hapiness == a.lastHapiness && a.hapiness >= 0 {
				a.iner = -1
				a.printf(ns, "decision: no need\n")
			} else if a.hapiness <= a.lastHapiness {
				a.fadeLastDecision(ns)
				if a.decided(ns) {
					if ns.log && ns.selected == a.id {
						a.printf(ns, "decision using network: %d outs: %s\n", a.direction, a.displayList(ns, a.outs))
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

func (a *Ant) decided(ns *Nests) bool {
	if ns.random {
		return false
	}
	if rand.Float64() < 0.1 {
		//return false
	}
	ins, ok := a.preparedEntries(a.entries)
	if !ok {
		a.printf(ns, "bad entries: %s\n", a.displayList(ns, a.entries))
		return false
	}
	a.outs = a.network.Propagate(ins, true)
	if ns.log {
		a.printf(ns, "Compute decision, propagation: %s\n", a.displayList(ns, a.outs))
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
		} else if a.dirCount < ns.bestAnt.dirCount-4 || (a.dirCount <= ns.bestAnt.dirCount && a.gRate < ns.bestAnt.gRate-20) {
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
		a.printf(ns, "network=%v hapiness=%.5f move: %d\n", a.networkDef, a.hapiness, a.direction)
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
	dist2Max := a.vision * a.vision
	dist2m := dist2Max
	var antMin *Ant
	a.contact = false
	for _, nest := range ns.nests {
		for _, ant := range nest.ants {
			if ant != a {
				dist2 := a.distAnt2(ant)
				if dist2 < dist2Max {
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
		}
	}
	if antMin != nil {
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
	}
	dist2m = dist2Max
	var foodMin *Food
	for _, food := range ns.foods {
		dist2 := a.distFood2(food)
		if dist2 < dist2Max {
			if dist2 < dist2Max/4 {
				a.statContact.incr()
				a.contact = true
			}
			if dist2 < dist2m {
				foodMin = food
				dist2m = dist2
			}
		}
	}
	if foodMin != nil {
		ang := math.Atan2(foodMin.X-a.x, foodMin.Y-a.y)
		if ang < 0 {
			ang = 2*math.Pi + ang
		}
		index := int(ang*float64(visionNb)/2.0/math.Pi + 0.000001)
		if index >= visionNb {
			index = index - visionNb
		}
		//a.printf(ns, "find %d angle=%0.2f degres=%0.2f index=%d\n", ant.id, ang, ang*180/math.Pi, index)
		a.entries[visionNb+index] = ((dist2Max - dist2m) / dist2Max)
	}
	a.computeHapiness(ns)
}

func (a *Ant) computeHapiness(ns *Nests) {
	if a.edge {
		a.hapiness = -.1
		return
	}
	a.lastHapiness = a.hapiness
	a.hapiness = 0
	for ii := visionNb; ii < visionNb*2; ii++ {
		a.hapiness += a.entries[ii]
	}
	if a.hapiness == 0 {
		for ii := 0; ii < visionNb; ii++ {
			a.hapiness -= a.entries[ii]
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
	if ns.random || a.lastDecision < 0 {
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
		a.printf(ns, "train %s => %v\n", a.displayList(ns, ins), a.outs)
		a.printf(ns, "outs: %s result=%f\n", a.displayList(ns, outs), trainResult)
	}
	if a.id == ns.selected {
		ns.addSample(ins, a.outs)
	}
	a.network.BackPropagate(a.outs)
	return true
}

func (a *Ant) fadeLastDecision(ns *Nests) bool {
	if ns.random || a.lastDecision == -1 {
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
		a.printf(ns, "fade %s => %v\n", a.displayList(ns, ins), a.outs)
		a.printf(ns, "outs: %s result=%.5f\n", a.displayList(ns, outs), trainResult)
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
