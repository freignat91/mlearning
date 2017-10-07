package network

import (
	"fmt"
	"math"
)

//MLNetwork .
type MLNetwork struct {
	netdef  []int
	in      *mlLayer
	layers  []*mlLayer
	out     *mlLayer
	phi     float64
	biais   bool
	created bool
}

var dataFileMap = make(map[string]*MlDataSet)

func sigmoid(value float64) float64 {
	return 1 / (1 + math.Exp(-value))
}

func sigmoidp(value float64) float64 {
	vv := sigmoid(value)
	return vv * (1 - vv)
}

// NewNetwork .
func NewNetwork(layers []int) (*MLNetwork, error) {
	if layers == nil || len(layers) < 2 {
		return nil, fmt.Errorf("Invalide layers definition")
	}
	n := &MLNetwork{
		netdef: layers,
		phi:    0.5,
		biais:  true,
	}
	n.in = n.newLayer("in", int(layers[0]))
	n.layers = make([]*mlLayer, len(layers)-2, len(layers)-2)
	preLayer := n.in
	for ii := 1; ii < len(layers)-1; ii++ {
		n.layers[ii-1] = n.newLayer(fmt.Sprintf("%d", ii), int(layers[ii]))
		preLayer.ConnectLayerTo(n.layers[ii-1])
		preLayer = n.layers[ii-1]
	}
	n.out = n.newLayer("out", int(layers[len(layers)-1]))
	preLayer.ConnectLayerTo(n.out)
	n.created = true
	return n, nil
}

// NetNetworkFromDataSet .
func NetNetworkFromDataSet(name string) (*MLNetwork, error) {
	dataSet, ok := dataFileMap[name]
	if !ok {
		return nil, fmt.Errorf("The logical train set %s doesn't exist", name)
	}
	return NewNetwork(dataSet.Layers)
}

// Getdef .
func (n *MLNetwork) Getdef() []int {
	return n.netdef
}

//IsCreated .
func (n *MLNetwork) IsCreated() bool {
	return n.created
}

// Display .
func (n *MLNetwork) Display(coef bool) []string {
	lines := make([]string, 0, 0)
	if n.in == nil {
		return lines
	}
	n.in.print(&lines)
	if coef {
		n.in.printCoef(&lines)
	}
	for _, layer := range n.layers {
		layer.print(&lines)
		if coef {
			layer.printCoef(&lines)
		}
	}
	n.out.print(&lines)
	return lines
}

//ClearIn .
func (n *MLNetwork) ClearIn() {
	for _, neuron := range n.in.neurons {
		neuron.value = 0
	}
}

//SetInValue .
func (n *MLNetwork) SetInValue(add bool, index int, value float64) {
	if add {
		n.in.neurons[index].value += value
	} else {
		n.in.neurons[index].value = value
	}
}

//GetInValue .
func (n *MLNetwork) GetInValue(index int) float64 {
	return n.in.neurons[index].value
}

// Propagate .
func (n *MLNetwork) Propagate(values []float64, retOut bool) []float64 {
	for ii, neuron := range n.in.neurons {
		if n.biais {
			if ii == 0 {
				neuron.value = 1
			} else {
				neuron.value = values[ii-1]
			}
		} else {
			neuron.value = values[ii]
		}
	}
	n.in.propagate()
	for _, layer := range n.layers {
		layer.propagate()
		//layer.print("layer")
	}
	if retOut {
		outs := n.getOutArray()
		//fmt.Printf("outs: %v\n", outs)
		return outs
	}
	return nil
}

func (n *MLNetwork) getOutArray() []float64 {
	ret := make([]float64, 0)
	for _, neuron := range n.out.neurons {
		ret = append(ret, neuron.value)
	}
	return ret
}

func (n *MLNetwork) matchRate(dataOut []float64) float64 {
	var errorRate float64
	for ii, neuron := range n.out.neurons {
		errorRate += ((dataOut[ii] - neuron.value) * (dataOut[ii] - neuron.value))
	}
	return math.Sqrt(errorRate)
}

// BackPropagate .
func (n *MLNetwork) BackPropagate(values []float64) {
	n.retroPropagateErrorDiff(values)
	n.updateCoef()
}

func (n *MLNetwork) retroPropagateErrorDiff(values []float64) {
	//compute out error diff
	for i, neuron := range n.out.neurons {
		neuron.errorDiff = (values[i] - neuron.value) // * sigmoidp(neuron.sum)
	}
	//compute all other layers error diff
	for l := len(n.layers) - 1; l >= 0; l-- {
		n.layers[l].retroPropagateErrorDiff()
	}
	n.in.retroPropagateErrorDiff()
}

func (n *MLNetwork) updateCoef() {
	n.in.updateCoef(n.phi)
	for _, layer := range n.layers {
		layer.updateCoef(n.phi)
	}
}

// Copy .
func (n *MLNetwork) Copy() (*MLNetwork, error) {
	net, err := NewNetwork(n.netdef)
	if err != nil {
		return nil, err
	}
	net.in.copyCoef(n.in)
	for ii, layer := range n.layers {
		net.layers[ii].copyCoef(layer)
	}
	net.out.copyCoef(n.out)
	return net, nil
}

//ComputeDistinctOut .
func (n *MLNetwork) ComputeDistinctOut() int {
	maxes := make(map[int]int)
	for ii := range n.in.neurons {
		ins := make([]float64, len(n.in.neurons), len(n.in.neurons))
		ins[ii] = 1
		outs := n.Propagate(ins, true)
		max := 0
		maxVal := 0.0
		for jj, val := range outs {
			if val > maxVal {
				maxVal = val
				max = jj
			}
		}
		maxes[max] = 1
	}
	return len(maxes)
}
