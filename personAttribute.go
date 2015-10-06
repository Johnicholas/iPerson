package main

type personAttribute struct {
	value float64
	// gameobject gameobject
}

func newPersonAttribute(startingValue float64) *personAttribute {
	answer := new(personAttribute)
	answer.value = startingValue
	answer.clamp()
	return answer
}

func (my *personAttribute) clamp() {
	if my.value > 1.0 {
		my.value = 1.0
	} else if my.value < 0.0 {
		my.value = 0.0
	} else {
		// it's ok
	}
}

func (my *personAttribute) boostAttribute(delta float64) {
	my.value += delta
	my.clamp()
}

func (my *personAttribute) update() {
	// do nothing?
}
