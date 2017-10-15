package nests

//Parameters .
type Parameters struct {
	workerInitNb                  int
	soldierInitNb                 int
	maxAnt                        int
	maxWorkerAnt                  int
	workerLife                    int
	soldierLife                   int
	workerMinSpeed                float64
	workerMaxSpeed                float64
	soldierMinSpeed               float64
	soldierMaxSpeed               float64
	workerLifeDecrPeriod          int64
	soldierLifeDecrPeriod         int64
	soldierInitCounter            int
	soldierRessourceCost          int
	pheromoneLevel                float64
	pheromoneAntDelay             int
	pheromoneGroup                float64
	pheromoneFadeDelay            int
	pheromoneFadeCounter          int
	nestAntRenewDelay             int64
	nestInitialRessource          int
	initialFoodGroupNumberPerNest int
	//
	networkUpdateDirCountDiff int
	networkUpdateGRateDiff    float64
}

func newParameters() *Parameters {
	return &Parameters{
		maxAnt:       800,
		maxWorkerAnt: 350,
		//
		workerInitNb:         50,
		workerLife:           120,
		workerMinSpeed:       0.3,
		workerMaxSpeed:       0.6,
		workerLifeDecrPeriod: 1000,
		//
		soldierInitNb:         0,
		soldierLife:           400,
		soldierMinSpeed:       0.25,
		soldierMaxSpeed:       0.6,
		soldierInitCounter:    3000,
		soldierRessourceCost:  4,
		soldierLifeDecrPeriod: 100,
		//
		pheromoneLevel:     1000,
		pheromoneAntDelay:  20,
		pheromoneFadeDelay: 3,
		pheromoneGroup:     100,
		//
		nestAntRenewDelay:             5000,
		nestInitialRessource:          300,
		initialFoodGroupNumberPerNest: 2,
		//
		networkUpdateDirCountDiff: 2,
		networkUpdateGRateDiff:    10,
	}
}
