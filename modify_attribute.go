package main

import "github.com/johnicholas/decisionflex"

type modifyAttribute struct {
	target     booster
	boostValue float64
}

func (my modifyAttribute) Perform(context decisionflex.Context) {
	my.target.boostAttribute(my.boostValue)
}
