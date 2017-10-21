package nests

//Parameters .
type Parameters struct {
	workerInitNb          int
	soldierInitNb         int
	maxAnt                int
	maxWorkerAnt          int
	workerLife            int
	soldierLife           int
	workerMinSpeed        float64
	workerMaxSpeed        float64
	soldierMinSpeed       float64
	soldierMaxSpeed       float64
	workerLifeDecrPeriod  int64
	soldierLifeDecrPeriod int64
	soldierInitCounter    int
	soldierRessourceCost  int
	pheromoneLevel        float64
	pheromoneAntDelay     int
	pheromoneGroup        float64
	pheromoneFadeDelay    int
	pheromoneFadeCounter  int
	//
	nestAntRenewDelay             int64
	nestInitialRessource          int
	initialFoodGroupNumberPerNest int
	initial2FoodGroupLenght       int
	initial4FoodGroupLenght       int
	//
	chanceToGetTheBestNetworkCopy float64 //[0-1]
	//
	updateTickNumber int64
}

func newParameters() *Parameters {
	return &Parameters{
		maxAnt:       800,
		maxWorkerAnt: 300,
		//
		workerInitNb:         100,
		workerLife:           120,
		workerMinSpeed:       0.6,
		workerMaxSpeed:       1,
		workerLifeDecrPeriod: 1000,
		//
		soldierInitNb:         0,
		soldierLife:           400,
		soldierMinSpeed:       0.5,
		soldierMaxSpeed:       0.9,
		soldierInitCounter:    3000,
		soldierRessourceCost:  4,
		soldierLifeDecrPeriod: 100,
		//
		pheromoneLevel:     1000,
		pheromoneAntDelay:  20,
		pheromoneFadeDelay: 5,
		pheromoneGroup:     100,
		//
		nestAntRenewDelay:             5000,
		nestInitialRessource:          500,
		initialFoodGroupNumberPerNest: 1,
		initial4FoodGroupLenght:       220,
		initial2FoodGroupLenght:       290,
		//
		chanceToGetTheBestNetworkCopy: 0.9,
		//
		updateTickNumber: 1000,
	}
}
