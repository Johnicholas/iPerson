package main

type modifyAttribute struct {
	target     booster
	boostValue float64
}

func (my modifyAttribute) Perform(context interface{}) {
	my.target.boostAttribute(my.boostValue)
}
