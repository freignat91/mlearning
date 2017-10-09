package nests

import "math/rand"

// FoodGroup .
type FoodGroup struct {
	y float64
	x float64
}

// Food .
type Food struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	carried   bool
	foodGroup *FoodGroup
}

func newFood(fg *FoodGroup) *Food {
	food := &Food{foodGroup: fg}
	food.renew()
	return food
}

func (f *Food) renew() {
	f.X = f.foodGroup.x + rand.Float64()*20
	f.Y = f.foodGroup.y + rand.Float64()*20
}
