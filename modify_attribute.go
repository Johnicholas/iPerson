// ModifyAttribute is a class, which is an Action
//
// ModifyAttribute implements Perform by
// calling BoostAttribute on a member, target, passing another member, boostValue

package main

import "github.com/johnicholas/decisionflex"

type modifyAttribute struct {
	decisionflex.Action
	target     booster
	boostValue float64
}

func (my *modifyAttribute) perform() {
	my.target.boostAttribute(my.boostValue)
}
