package main

import "fmt"
import "github.com/johnicholas/decisionflex"

const (
	EJECT    = iota
	SLIDE    = iota
	TIP      = iota
	SAMPLE   = iota
	DISPENSE = iota
	SHUCKER  = iota
)

const (
	SLIDE_NONE = iota
	DRY        = iota
	WET        = iota
)

const (
	TIP_NONE = iota
	CLEAN    = iota
	FULL     = iota
	DIRTY    = iota
)

type catContext struct {
	DryScheduled int
	WetScheduled int
	RobotAt      int // TODO: one of EJECT, SLIDE, etc
	SlideIs      int // TODO: one of DRY, WET
	TipIs        int // TODO: one of CLEAN, FULL, etc
}

type catContextMutator func(*catContext)

// catContext mutator methods can be Performers,
// by casting the incoming context to a catContext
func (self catContextMutator) Perform(context interface{}) {
	self(context.(*catContext))
}

type goTo struct {
	Destination int
}

func (g goTo) Perform(context interface{}) {
	context.(*catContext).RobotAt = g.Destination
}

// TODO: catContext accessor methods can be Considerers,
// by casting the incoming context to a catContext?

type catContextAccessor func(catContext) bool

func (self catContextAccessor) Consider(context interface{}) float64 {
	if self(*(context.(*catContext))) {
		return 1.0
	} else {
		return 0.0
	}
}

var someDryScheduled catContextAccessor = func(c catContext) bool { return c.DryScheduled > 0 }
var noDryScheduled catContextAccessor = func(c catContext) bool { return c.DryScheduled == 0 }
var someWetScheduled catContextAccessor = func(c catContext) bool { return c.WetScheduled > 0 }
var noWetScheduled catContextAccessor = func(c catContext) bool { return c.WetScheduled == 0 }

type robotAt struct {
	Location int // TODO: one of EJECT, SLIDE, etc
}

func (r robotAt) Consider(context interface{}) float64 {
	if context.(*catContext).RobotAt == r.Location {
		return 1.0
	} else {
		return 0.0
	}
}

type slideIs struct {
	State int // TODO: one of DRY, WET, etc
}

func (s slideIs) Consider(context interface{}) float64 {
	if context.(*catContext).SlideIs == s.State {
		return 1.0
	} else {
		return 0.0
	}
}

type tipIs struct {
	State int // TODO: one of TIP_NONE, CLEAN, FULL, etc
}

func (t tipIs) Consider(context interface{}) float64 {
	if context.(*catContext).TipIs == t.State {
		return 1.0
	} else {
		return 0.0
	}
}

type tipIsNot struct {
	State int // TODO: one of TIP_NONE, CLEAN, FULL, etc
}

func (t tipIsNot) Consider(context interface{}) float64 {
	if context.(*catContext).TipIs != t.State {
		return 1.0
	} else {
		return 0.0
	}
}

type firstPossibleT struct{}

func (c firstPossibleT) Choose(choices []decisionflex.ActionSelection) decisionflex.ActionSelection {
	for _, choice := range choices {
		if choice.Score > 0.0 {
			return choice
		}
	}
	return choices[len(choices)-1]
}

func (c firstPossibleT) String() string {
	return "choose first possible"
}

var firstPossible firstPossibleT

func main() {
	context := catContext{
		DryScheduled: 3,
		WetScheduled: 21,
		RobotAt:      SLIDE,
		SlideIs:      DRY,
		// TipIs: NONE
	}

	// if (dry>0&&dry_scheduled>0) fire(acquire_dry)
	possiblyAcquireDry := decisionflex.NewActionConsiderations("acquire a dry reading")
	possiblyAcquireDry.AddConsiderer(slideIs{DRY})
	possiblyAcquireDry.AddConsiderer(someDryScheduled)
	possiblyAcquireDry.AddPerformer(catContextMutator(func(self *catContext) {
		// preconditions
		if self.DryScheduled <= 0 {
			panic("unscheduled acquire dry")
		}
		// effects
		self.DryScheduled--
	}))

	// if (dry>0&&dry_scheduled==0&&at_slide>0&&clean==0) fire(at_tip)
	possiblyGoToTip := decisionflex.NewActionConsiderations("go to tip load station")
	possiblyGoToTip.AddConsiderer(slideIs{DRY})
	possiblyGoToTip.AddConsiderer(noDryScheduled)
	possiblyGoToTip.AddConsiderer(robotAt{SLIDE})
	possiblyGoToTip.AddConsiderer(tipIsNot{CLEAN})
	possiblyGoToTip.AddPerformer(goTo{TIP})

	// if (dry>0&&dry_scheduled==0&&clean==0&&full==0&&at_tip>0) fire(load_tip)
	possiblyLoadTip := decisionflex.NewActionConsiderations("load a tip")
	possiblyLoadTip.AddConsiderer(slideIs{DRY})
	possiblyLoadTip.AddConsiderer(noDryScheduled)
	possiblyLoadTip.AddConsiderer(tipIsNot{CLEAN})
	possiblyLoadTip.AddConsiderer(tipIsNot{FULL})
	possiblyLoadTip.AddConsiderer(robotAt{TIP})
	possiblyLoadTip.AddPerformer(catContextMutator(func(self *catContext) {
		// preconditions
		if self.RobotAt != TIP {
			panic("load tip when robot is not at tip load station")
		}
		if self.TipIs != TIP_NONE {
			panic("load tip with another tip already on proboscis")
		}
		// effects
		self.TipIs = CLEAN
	}))

	// if(clean>0&&full==0&&at_tip>0) fire(at_slide)
	//possiblyGoToSlide := decisionflex.NewActionConsiderations("go to slide load station")
	//possiblyGoToSlide.AddConsiderer(tipIs{CLEAN})
	//possiblyGoToSlide.AddConsiderer(tipIsNot{FULL})
	//possiblyGoToSlide.AddConsiderer(robotAt{TIP})
	//possiblyGoToSlide.AddPerformer(goTo{SLIDE})

	// if (clean>0&&full==0&&at_slide>0) fire(at_sample)
	possiblyGoToSample := decisionflex.NewActionConsiderations("go to sample cup")
	possiblyGoToSample.AddConsiderer(tipIs{CLEAN})
	possiblyGoToSample.AddConsiderer(tipIsNot{FULL})
	possiblyGoToSample.AddConsiderer(robotAt{TIP})
	possiblyGoToSample.AddPerformer(goTo{SAMPLE})

	// if (clean>0&&full==0&&at_sample>0) fire(aspirate)
	possiblyAspirate := decisionflex.NewActionConsiderations("aspirate from sample cup")
	possiblyAspirate.AddConsiderer(tipIs{CLEAN})
	possiblyAspirate.AddConsiderer(tipIsNot{FULL})
	possiblyAspirate.AddConsiderer(robotAt{SAMPLE})
	possiblyAspirate.AddPerformer(catContextMutator(func(self *catContext) {
		// preconditions
		if self.RobotAt != SAMPLE {
			panic("aspirate while robot is not at sample cup")
		}
		if self.TipIs != CLEAN {
			panic("aspirate while tip is not clean")
		}
		// effects
		self.TipIs = FULL
	}))

	// if (dry>0&&full>0&&at_sample>0) fire(at_dispense)
	possiblyGoToDispense := decisionflex.NewActionConsiderations("go to dispense station")
	possiblyGoToDispense.AddConsiderer(slideIs{DRY})
	possiblyGoToDispense.AddConsiderer(tipIs{FULL})
	possiblyGoToDispense.AddConsiderer(robotAt{SAMPLE})
	possiblyGoToDispense.AddPerformer(goTo{DISPENSE})

	// if (dry>0&&full>0&&at_dispense) fire(dispense_on)
	possiblyDispense := decisionflex.NewActionConsiderations("dispense sample onto slide")
	possiblyDispense.AddConsiderer(slideIs{DRY})
	possiblyDispense.AddConsiderer(tipIs{FULL})
	possiblyDispense.AddConsiderer(robotAt{DISPENSE})
	possiblyDispense.AddPerformer(catContextMutator(func(self *catContext) {
		// preconditions
		if self.SlideIs != DRY {
			panic("dispense on slide when slide is not dry")
		}
		if self.TipIs != FULL {
			panic("dispense on slide when tip is not full")
		}
		// effects
		self.SlideIs = WET
		self.TipIs = DIRTY
	}))

	// if (dirty>0&&at_dispense>0) fire(at_shucker)
	possiblyGoToShucker := decisionflex.NewActionConsiderations("go to tip shucker")
	possiblyGoToShucker.AddConsiderer(tipIs{DIRTY})
	possiblyGoToShucker.AddConsiderer(robotAt{DISPENSE})
	possiblyGoToShucker.AddPerformer(goTo{SHUCKER})

	// if (dirty>0&&at_shucker>0) fire(shuck_tip)
	possiblyShuckTip := decisionflex.NewActionConsiderations("shuck the tip")
	possiblyShuckTip.AddConsiderer(tipIs{DIRTY})
	possiblyShuckTip.AddConsiderer(robotAt{SHUCKER})
	possiblyShuckTip.AddPerformer(catContextMutator(func(self *catContext) {
		// preconditions
		if self.RobotAt != SHUCKER {
			panic("shuck tip while robot is not at tip shucker")
		}
		// effects
		self.TipIs = TIP_NONE
	}))

	// if (wet>0&&wet_schedule>0) fire(acquire_wet)
	possiblyAcquireWet := decisionflex.NewActionConsiderations("acquire a wet reading")
	possiblyAcquireWet.AddConsiderer(slideIs{WET})
	possiblyAcquireWet.AddConsiderer(someWetScheduled)
	possiblyAcquireWet.AddPerformer(catContextMutator(func(self *catContext) {
		// preconditions
		if self.WetScheduled <= 0 {
			panic("unscheduled acquire wet")
		}
		// effects
		self.WetScheduled--
	}))

	// if (wet>0&&at_shucker>0) fire(at_eject)
	possiblyGoToEject := decisionflex.NewActionConsiderations("go to eject station")
	possiblyGoToEject.AddConsiderer(slideIs{WET})
	possiblyGoToEject.AddConsiderer(robotAt{SHUCKER})
	possiblyGoToEject.AddPerformer(goTo{EJECT})

	// if (wet>0&&at_eject>0) fire(eject)
	possiblyEject := decisionflex.NewActionConsiderations("eject a slide")
	possiblyEject.AddConsiderer(slideIs{WET})
	possiblyEject.AddConsiderer(robotAt{EJECT})
	possiblyEject.AddPerformer(catContextMutator(func(self *catContext) {
		// preconditions
		if self.RobotAt != EJECT {
			panic("eject when robot is not at eject station")
		}
		if self.TipIs != TIP_NONE {
			panic("eject when proboscis has a tip on it")
		}
		// effect
		self.SlideIs = SLIDE_NONE
	}))

	idle := decisionflex.NewActionConsiderations("nothing to do!")

	decider := decisionflex.New(
		decisionflex.SingleContextFactory{&context},
		firstPossible,
	)
	decider.Add(possiblyAcquireDry)
	decider.Add(possiblyGoToTip)
	decider.Add(possiblyLoadTip)
	// decider.Add(possiblyGoToSlide)
	decider.Add(possiblyGoToSample)
	decider.Add(possiblyAspirate)
	decider.Add(possiblyGoToDispense)
	decider.Add(possiblyDispense)
	decider.Add(possiblyGoToShucker)
	decider.Add(possiblyShuckTip)
	decider.Add(possiblyAcquireWet)
	decider.Add(possiblyGoToEject)
	decider.Add(possiblyEject)
	decider.Add(idle)

	for i := 0; i < 100; i++ {
		answer := decider.PerformAction()
		if answer.ActionObject == idle.ActionObject {
			break
		} else {
			fmt.Println(answer.ActionObject)
		}

	}
}
