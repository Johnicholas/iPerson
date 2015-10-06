package main

import "github.com/johnicholas/decisionflex"
import "math/rand"

func uniform(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

type randomModifyAttribute struct {
	target   booster
	min, max float64
}

func (my randomModifyAttribute) Perform(context decisionflex.Context) {
	my.target.boostAttribute(uniform(my.min, my.max))

}
