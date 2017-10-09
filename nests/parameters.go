package nests

//Parameters .
type Parameters struct {
	pheromoneLevel       float64
	pheromoneAntDelay    int
	pheromoneGroup       float64
	pheromoneFadeDelay   int
	pheromoneFadeCounter int
}

func newParameters() *Parameters {
	return &Parameters{
		pheromoneLevel:     1000,
		pheromoneAntDelay:  20,
		pheromoneFadeDelay: 150,
		pheromoneGroup:     100,
	}
}
