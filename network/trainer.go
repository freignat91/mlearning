package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

//MlDataSet .
type MlDataSet struct {
	Name    string         `json:"name"`
	Layers  []int          `json:"layers"`
	Data    []MlDataSample `json:"data"`
	nbTrain int
}

//MlDataSample .
type MlDataSample struct {
	In  []float64 `json:"in"`
	Out []float64 `json:"out"`
}

// LoadTrainFile .
func (n *MLNetwork) LoadTrainFile(path string) ([]string, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var dataSet MlDataSet
	json.Unmarshal(raw, &dataSet)
	if dataSet.Name == "" {
		return nil, fmt.Errorf("dataSet file do not have name")
	}
	dataFileMap[dataSet.Name] = &dataSet
	lines := make([]string, 0, 0)
	lines = append(lines, fmt.Sprintf("File %s loaded, sample number: %d\n", dataSet.Name, len(dataSet.Data)))
	lines = append(lines, fmt.Sprintf("layers %v\n", dataSet.Layers))
	return lines, nil
}

//Train .
func (n *MLNetwork) Train(name string, nbl int, all bool, hide bool, createNetwork bool, analyse bool) ([]string, error) {
	dataSet, ok := dataFileMap[name]
	if !ok {
		return nil, fmt.Errorf("The logical train set %s doesn't exist", name)
	}
	lines := make([]string, 0, 0)
	if analyse {
		n.analyseDataSet(&lines, dataSet)
	}
	start := time.Now()
	nb := nbl
	if all {
		nb = nb * len(dataSet.Data)
	}
	lines = append(lines, fmt.Sprintf("Train network with %s for %d iterations\n", dataSet.Name, nb))
	for ii := 0; ii < nb; ii++ {
		data := dataSet.Data[rand.Int31n(int32(len(dataSet.Data)))]
		n.Propagate(data.In, false)
		n.BackPropagate(data.Out)
		dataSet.nbTrain++
		if time.Now().Sub(start) > time.Second*10 {
			n.trainResult(&lines, name, ii, hide, dataSet)
			start = time.Now()
		}
	}
	n.trainResult(&lines, name, nb, hide, dataSet)
	return lines, nil
}

func (n *MLNetwork) trainResult(lines *[]string, name string, nb int, hide bool, dataSet *MlDataSet) {
	*lines = append(*lines, fmt.Sprintf("Result for dataSet %s at %d:\n", name, dataSet.nbTrain))
	var matchRate float64
	tot := 0.0
	maxes := make(map[int]int)
	for _, data := range dataSet.Data {
		n.Propagate(data.In, false)
		matchRate += n.matchRate(data.Out)
		if !hide {
			line, val, max := n.sList(n.getOutArray(), "%.2f", true, -1)
			maxes[max] = 1
			tot += val
			*lines = append(*lines, fmt.Sprintf("%v => %v: %s\n", data.In, data.Out, line))
		}
	}
	matchRate = matchRate / float64(len(dataSet.Data))
	*lines = append(*lines, fmt.Sprintf("after %d iteration, match rate:%.5f tot:%.5f distinct=%d\n", nb, matchRate, tot, len(maxes)))
}

func (n *MLNetwork) sList(list []float64, format string, all bool, braq int) (string, float64, int) {
	moy := 0.0
	max := 0.0
	maxii := 0
	ret := "[ "
	braqii := 0
	for ii, val := range list {
		if val > max {
			max = val
			maxii = ii
		}
		moy += val
		if ii == braq {
			braqii = ii
			ret += fmt.Sprintf("("+format+") ", val)
		} else {
			ret += fmt.Sprintf(format+" ", val)
		}
	}
	moy = moy / float64(len(list))
	if all {
		if braq >= 0 {
			return ret + fmt.Sprintf("] max=%d/%d diffMax=%.2f", maxii, braqii, max-moy), max - moy, maxii
		}
		return ret + fmt.Sprintf("] max=%d diffMax=%.2f", maxii, max-moy), max - moy, maxii
	}
	return ret + "]", max - moy, maxii
}

func (n MLNetwork) analyseDataSet(lines *[]string, dataSet *MlDataSet) {
	nin := dataSet.Layers[0]
	nout := dataSet.Layers[len(dataSet.Layers)-1]

	insd := make([][]float64, nin, nin)
	for ii := range insd {
		insd[ii] = make([]float64, nout, nout)
	}
	ins := make([]float64, nin, nin)
	outs := make([]float64, nout, nout)
	for _, data := range dataSet.Data {
		for ii, val := range data.In {
			ins[ii] += val
			if val > 0 {
				inArray := insd[ii]
				for jj, valo := range data.Out {
					inArray[jj] += valo
				}
			}
		}
		for ii, val := range data.Out {
			outs[ii] += val
		}
	}
	*lines = append(*lines, fmt.Sprintln("analyse:"))
	*lines = append(*lines, fmt.Sprintf("moy in: %v\n", ins))
	*lines = append(*lines, fmt.Sprintf("moy ou: %v\n", outs))
	for ii := range insd {
		sol := ii + nout/2
		if sol >= nout {
			sol = sol - nout
		}
		sret := "[ "
		for jj, val := range insd[ii] {
			if jj == sol {
				sret += fmt.Sprintf("(%.0f) ", val)
			} else {
				sret += fmt.Sprintf("%.0f ", val)
			}
		}
		sret += "]"
		*lines = append(*lines, fmt.Sprintf("in %d: %s\n", ii, sret))
	}
}
