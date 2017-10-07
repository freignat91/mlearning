package nests

// Stats .
type Stats struct {
	value  int64
	cumul  int64
	scumul int64
	nest   *Stats
	nests  *Stats
}

func newStats(ns *Stats, n *Stats) *Stats {
	return &Stats{
		nests: ns,
		nest:  n,
	}
}

func (s *Stats) incr() {
	s.value++
	s.scumul++
	if s.nest != nil {
		s.nest.incr()
	}
	if s.nests != nil {
		s.nests.incr()
	}
}

func (s *Stats) push() {
	s.cumul = s.value
	s.value = 0
}
