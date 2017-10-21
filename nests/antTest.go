package nests

import "fmt"

func (a *Ant) trainSoluce(ns *Nests, nb int) {
	a.trained = true
	//a.Life = 100000000
	ins := []int{0, 1, 2, 3}
	if a.AntType == 1 {
		ins = []int{0, 3}
	}
	for _, in := range ins {
		direct := 0
		for ii := 0; ii < nb; ii++ {
			a.setEntriesSoluce(in, direct)
			a.setOutsSoluce(in, direct)
			a.network.Propagate(a.entries, false)
			a.network.BackPropagate(a.outs)
			direct++
			if direct >= outNb {
				direct = 0
			}
		}
	}
}

func (a *Ant) setEntriesSoluce(in int, direct int) {
	index := in * visionNb
	for ii := range a.entries {
		if ii == index+direct {
			a.entries[ii] = 1
		} else {
			a.entries[ii] = 0
		}
	}
}

func (a *Ant) setOutsSoluce(in int, direct int) {
	if in == 0 || (a.AntType == 0 && in == 3) {
		direct = direct + outNb/2
		if direct >= outNb {
			direct = direct - outNb
		}
	}
	for ii := range a.outs {
		if ii == direct {
			a.outs[ii] = 1
		} else {
			a.outs[ii] = 0
		}
	}
}

func (a *Ant) test(ns *Nests) []string {
	lines := make([]string, 1, 1)
	lines[0] = fmt.Sprintf("Test selected network: %v\n", a.network.Getdef())
	ins := []int{0, 1, 2, 3}
	if a.AntType == 1 {
		ins = []int{0, 3}
	}
	for _, in := range ins {
		dirMap := make(map[int]int)
		for direct := 0; direct < 8; direct++ {
			a.setEntriesSoluce(in, direct)
			outs := a.network.Propagate(a.entries, true)
			max := getMax(outs)
			dirMap[max] = 1
			lines = append(lines, fmt.Sprintf("in: %s out=%s max=%d\n", a.displayList(ns, a.entries, "%.0f"), a.displayList(ns, outs, "%.3f"), max))
		}
		lines = append(lines, fmt.Sprintf("Test for entries=%d distinct=%d\n", in, len(dirMap)))
	}
	return lines
}

func getMax(list []float64) int {
	max := 0
	maxVal := 0.0
	for ii, val := range list {
		if val > maxVal {
			maxVal = val
			max = ii
		}
	}
	return max
}
