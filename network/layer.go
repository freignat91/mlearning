package network

import (
	"fmt"
	"math/rand"
)

type mlLayer struct {
	id        string
	neurons   []*mlNeuron
	linksTo   []*mlLink
	nextLayer *mlLayer
	isIn      bool
	isOut     bool
	biais     bool
}

type mlNeuron struct {
	id        string
	sum       float64
	value     float64
	errorDiff float64
}

func newNeuron(layerID string, id int) *mlNeuron {
	return &mlNeuron{
		id: fmt.Sprintf("%s-%d", layerID, id),
	}
}

func (n *MLNetwork) newLayer(id string, size int) *mlLayer {
	layer := mlLayer{
		id:    id,
		biais: n.biais,
	}
	if id == "in" {
		layer.isIn = true
	}
	if id == "out" {
		layer.isOut = true
	}
	if n.biais && !layer.isOut && !layer.isIn {
		size++
	}
	layer.neurons = make([]*mlNeuron, size, size)
	layer.linksTo = make([]*mlLink, 0)
	for ii := range layer.neurons {
		layer.neurons[ii] = newNeuron(layer.id, ii)
	}
	if n.biais && !layer.isOut {
		layer.neurons[0].value = 1
	}
	return &layer
}

func (l *mlLayer) print(lines *[]string) {
	*lines = append(*lines, fmt.Sprintf("Layer %s:\n", l.id))
	for _, neuron := range l.neurons {
		*lines = append(*lines, fmt.Sprintf("neuron %+v\n", neuron))
	}
}

func (l *mlLayer) printCoef(lines *[]string) {
	ret := "["
	for _, link := range l.linksTo {
		ret += fmt.Sprintf("%s<->%s: %.3f, ", link.neuronFrom.id, link.neuronTo.id, link.coef)
	}
	*lines = append(*lines, fmt.Sprintf("Coef: %s]\n", ret))
}

func (l *mlLayer) ConnectLayerTo(layer *mlLayer) {
	//fmt.Printf("connect layer=%s (%d) to %s (%d)\n", l.id, len(l.neurons), layer.id, len(layer.neurons))
	for _, neuronFrom := range l.neurons {
		for ii, neuronTo := range layer.neurons {
			if layer.isOut || ii > 0 || !l.biais {
				link := mlLink{
					coef:       rand.Float64(),
					neuronFrom: neuronFrom,
					neuronTo:   neuronTo,
				}
				l.linksTo = append(l.linksTo, &link)
			}
		}
	}
	l.nextLayer = layer
}

func (l *mlLayer) propagate() {
	for _, neuron := range l.nextLayer.neurons {
		neuron.sum = 0
	}
	for _, link := range l.linksTo {
		link.propagate()
	}
	for _, neuron := range l.nextLayer.neurons {
		neuron.value = sigmoid(neuron.sum)
	}
}

func (l *mlLayer) retroPropagateErrorDiff() {
	for _, neuron := range l.neurons {
		neuron.errorDiff = 0
	}
	for _, link := range l.linksTo {
		link.retroPropagateErrorDiff()
	}
}

func (l *mlLayer) updateCoef(phi float64) {
	for _, link := range l.linksTo {
		link.updateCoef(phi)
	}
}

func (l *mlLayer) copyCoef(layer *mlLayer) {
	for ii, link := range layer.linksTo {
		l.linksTo[ii].coef = link.coef
	}
}
