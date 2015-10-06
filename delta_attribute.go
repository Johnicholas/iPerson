package main

// import "time" // TODO?

type deltaAttribute struct {
	personAttribute
	delta_per_second float64
}

func (my *deltaAttribute) update() {
	my.personAttribute.update()
	my.value += my.delta_per_second // delta_per_second * some_elapsed_time()?
}
