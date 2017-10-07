package network

type mlLink struct {
	coef        float64
	updatedCoef float64
	neuronFrom  *mlNeuron
	neuronTo    *mlNeuron
}

func (l *mlLink) propagate() {
	l.neuronTo.sum += l.coef * l.neuronFrom.value
}

//lyser.linksTo
func (l *mlLink) retroPropagateErrorDiff() {
	l.neuronFrom.errorDiff += l.coef * l.neuronTo.errorDiff
}

//layer.linkFrom
func (l *mlLink) updateCoef(phi float64) {
	l.coef += phi * l.neuronTo.errorDiff * sigmoidp(l.neuronTo.sum) * l.neuronFrom.value
	//fmt.Printf("%s->%s: errorDiff=%f sigmoidp(sum)=%f value=%f result=%f\n", l.neuronFrom.id, l.neuronTo.id, l.neuronTo.errorDiff, sigmoidp(l.neuronTo.sum), l.neuronFrom.value, l.coef)
}
