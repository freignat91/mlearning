package nests

//Nest .
type Nest struct {
	id            int
	ants          []*Ant
	statTrain     *Stats
	statDecision  *Stats
	statReinforce *Stats
	statFade      *Stats
	statNetwork   *Stats
	statContact   *Stats
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

func (n *Nest) addData(list *[]*Data) {
	for _, ant := range n.ants {
		*list = append(*list, ant.getData())
	}
}

func (n *Nest) nextTime(ns *Nests) {
	h := 0.0
	for _, ant := range n.ants {
		ant.nextTime(ns)
		h += ant.hapiness
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
