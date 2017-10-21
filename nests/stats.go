package nests

// Stats .
type Stats struct {
	value  int64
	cumul  int64
	scumul int64
	nest   *Stats
}

func newStats(n *Stats) *Stats {
	return &Stats{
		nest: n,
	}
}

func (s *Stats) incr() {
	s.value++
	s.scumul++
	if s.nest != nil {
		s.nest.incr()
	}
}

func (s *Stats) push() {
	s.cumul = s.value
	s.value = 0
}
