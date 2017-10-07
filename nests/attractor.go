package nests

const aMeetSame = "meet"

//Attractors .
type Attractors struct {
	nameMap map[string]bool
}

func newAttractors() *Attractors {
	att := &Attractors{
		nameMap: make(map[string]bool),
	}
	att.nameMap[aMeetSame] = false
	return att
}
