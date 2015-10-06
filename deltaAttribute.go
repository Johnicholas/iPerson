// DeltaAttribute is a class, which is a PersonAttribute
//
// DeltaAttribute implements Update by calling base.Update,
// and also increasing Value by delta_per_second (a member) multiplied by deltaTime
// deltaTime is apparently accessible from the Time module?

package main

// import "time"

type deltaAttribute struct {
	personAttribute
	delta_per_second float64
}

func (my *deltaAttribute) update() {
	my.personAttribute.update()
	my.value += my.delta_per_second // delta_per_second * some_elapsed_time()?
}
