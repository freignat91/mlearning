package nests

import "math/rand"

// FoodGroup .
type FoodGroup struct {
	y float64
	x float64
}

// Food .
type Food struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	carried bool
}

func newFood(fg *FoodGroup) *Food {
	food := &Food{}
	fg.setPosition(food)
	return food
}

func (fg *FoodGroup) setPosition(f *Food) {
	f.X = fg.x + rand.Float64()*20
	f.Y = fg.y + rand.Float64()*20
}
