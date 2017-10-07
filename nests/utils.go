package nests

import "fmt"

func (a *Ant) getCloser(ns *Nests) []*Ant {
	distMax := a.vision * a.vision
	list := make([]*Ant, 0, 0)
	for _, nest := range ns.nests {
		for _, ant := range nest.ants {
			if a.dist2(ant) < distMax {
				list = append(list, ant)
			}
		}
	}
	return list
}

func (a *Ant) dist2(ant *Ant) float64 {
	return (ant.x-a.x)*(ant.x-a.x) + (ant.y-a.y)*(ant.y-a.y)
}

func (a *Ant) printf(ns *Nests, format string, args ...interface{}) {
	if ns.log && a.id == ns.selected {
		fmt.Printf(format, args...)
	}
}

func (a *Ant) displayList(ns *Nests, list []float64) string {
	if a.id != ns.selected {
		return ""
	}
	ret := "[ "
	for _, val := range list {
		ret += fmt.Sprintf("%.3f ", val)
	}
	return ret + "]"
}

func (a *Ant) computeTrainResult(ref []float64, outs []float64) float64 {
	var ret float64
	for ii, out := range outs {
		ret += ((ref[ii] - out) * (ref[ii] - out))
	}
	return ret
}
