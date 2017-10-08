package nests

import "math/rand"

//Nest .
type Nest struct {
	id            int
	x             float64
	y             float64
	ants          []*Ant
	statTrain     *Stats
	statDecision  *Stats
	statReinforce *Stats
	statFade      *Stats
	statNetwork   *Stats
	statContact   *Stats
	statFood      *Stats
	bestAnt       *Ant
	worseAnt      *Ant
	happiness     float64
}

func newNest(ns *Nests, id int, nb int) (*Nest, error) {
	nest := &Nest{
		id:            id,
		ants:          make([]*Ant, nb, nb),
		statTrain:     newStats(nil, nil),
		statDecision:  newStats(nil, nil),
		statReinforce: newStats(nil, nil),
		statFade:      newStats(nil, nil),
		statNetwork:   newStats(nil, nil),
		statContact:   newStats(nil, nil),
		statFood:      newNestStats(ns.statFood),
	}
	if rand.Float64() < 0.5 {
		nest.x = ns.xmin + 25
		nest.y = ns.ymin + 25
	} else {
		nest.x = ns.xmax - 25
		nest.y = ns.ymax - 25
	}
	for ii := range nest.ants {
		ant, err := newAnt(ns, nest, ii+1)
		if err != nil {
			return nil, err
		}
		nest.ants[ii] = ant
	}
	nest.bestAnt = nest.ants[0]
	nest.worseAnt = nest.ants[0]
	return nest, nil
}

func (n *Nest) addData(list *[]*AntData) {
	for _, ant := range n.ants {
		*list = append(*list, ant.getData())
	}
}

func (n *Nest) nextTime(ns *Nests) {
	h := 0.0
	for _, ant := range n.ants {
		ant.nextTime(ns)
		h += ant.happiness
		if ant.dirCount > n.bestAnt.dirCount {
			n.bestAnt = ant
		} else if ant.dirCount == n.bestAnt.dirCount && ant.gRate > n.bestAnt.gRate {
			n.bestAnt = ant
		}
		if ant.dirCount < n.worseAnt.dirCount {
			n.worseAnt = ant
		} else if ant.dirCount == n.worseAnt.dirCount && ant.gRate < n.worseAnt.gRate {
			n.worseAnt = ant
		}
	}
	n.happiness = h / float64(len(n.ants))
}
