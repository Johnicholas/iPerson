package main

import (
	"flag"
	"fmt"
	"github.com/johnicholas/decisionflex"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Location int

const (
	EJECT Location = iota
	SLIDE
	TIP
	SAMPLE
	DISPENSE
	SHUCKER
)

type SlideState int

const (
	NO_SLIDE SlideState = iota
	DRY
	WET
)

type TipState int

const (
	NO_TIP TipState = iota
	CLEAN
	FULL
	DIRTY
)

type Context struct {
	DryScheduled int
	WetScheduled int
	RobotAt      Location
	SlideIs      SlideState
	TipIs        TipState
}

// TODO: isn't there some reflection-ish way to do this?
func (c Context) String() string {
	s := []string{
		fmt.Sprintf("{\n"),
		fmt.Sprintf("\tDryScheduled: %d\n", c.DryScheduled),
		fmt.Sprintf("\tWetScheduled: %d\n", c.WetScheduled),
		fmt.Sprintf("\tRobotAt: %d\n", c.RobotAt),
		fmt.Sprintf("\tSlideIs: %d\n", c.SlideIs),
		fmt.Sprintf("\tTipIs: %d\n", c.TipIs),
		fmt.Sprintf("}\n"),
	}
	return strings.Join(s, "")
}

type catPrecondition struct {
	test func(Context) bool
}

func (c catPrecondition) Consider(i interface{}) float64 {
	downcast := i.(*Context)
	if c.test(*downcast) {
		return 1.0
	} else {
		return 0.0
	}
}

var someDryScheduled = catPrecondition{func(c Context) bool {
	return c.DryScheduled > 0
}}

var noDryScheduled = catPrecondition{func(c Context) bool {
	return c.DryScheduled == 0
}}

var someWetScheduled = catPrecondition{func(c Context) bool {
	return c.WetScheduled > 0
}}

var noWetScheduled = catPrecondition{func(c Context) bool {
	return c.WetScheduled == 0
}}

func robotAt(l Location) catPrecondition {
	return catPrecondition{func(c Context) bool {
		return c.RobotAt == l
	}}
}

func robotNotAt(l Location) catPrecondition {
	return catPrecondition{func(c Context) bool {
		return c.RobotAt != l
	}}
}

func slideIs(s SlideState) catPrecondition {
	return catPrecondition{func(c Context) bool {
		return c.SlideIs == s
	}}
}

func tipIs(t TipState) catPrecondition {
	return catPrecondition{func(c Context) bool {
		return c.TipIs == t
	}}
}

func tipIsNot(t TipState) catPrecondition {
	return catPrecondition{func(c Context) bool {
		return c.TipIs != t
	}}
}

type catMutator struct {
	toRun func(*Context)
}

func (c catMutator) Perform(i interface{}) {
	downcast := i.(*Context)
	c.toRun(downcast)
}

func goTo(l Location) catMutator {
	return catMutator{func(c *Context) {
		c.RobotAt = l
	}}
}

func main() {
	context := Context{
		DryScheduled: 3,
		WetScheduled: 21,
		RobotAt:      SLIDE,
		SlideIs:      DRY,
		// TipIs: NO_TIP
	}

	goToTip := decisionflex.Rule{
		Name: "go to tip load station",
		Considerations: []decisionflex.Consideration{
			slideIs(DRY),
			tipIs(NO_TIP),
			robotNotAt(TIP),
		},
		Actions: []decisionflex.Action{goTo(TIP)},
	}

	loadTip := decisionflex.Rule{
		Name: "load a tip",
		Considerations: []decisionflex.Consideration{
			slideIs(DRY),
			tipIs(NO_TIP),
			robotAt(TIP),
		},
		Actions: []decisionflex.Action{catMutator{func(c *Context) {
			c.TipIs = CLEAN
		}}},
	}

	goToSample := decisionflex.Rule{
		Name: "go to sample cup",
		Considerations: []decisionflex.Consideration{
			slideIs(DRY),
			tipIs(CLEAN),
			robotNotAt(SAMPLE),
		},
		Actions: []decisionflex.Action{goTo(SAMPLE)},
	}

	aspirate := decisionflex.Rule{
		Name: "aspirate from sample cup",
		Considerations: []decisionflex.Consideration{
			slideIs(DRY),
			noDryScheduled,
			tipIs(CLEAN),
			robotAt(SAMPLE),
		},
		Actions: []decisionflex.Action{catMutator{func(c *Context) {
			c.TipIs = FULL
		}}},
	}

	goToDispense := decisionflex.Rule{
		Name: "go to dispense station",
		Considerations: []decisionflex.Consideration{
			slideIs(DRY),
			noDryScheduled,
			tipIs(FULL),
			robotNotAt(DISPENSE),
		},
		Actions: []decisionflex.Action{goTo(DISPENSE)},
	}

	dispense := decisionflex.Rule{
		Name: "dispense sample onto slide",
		Considerations: []decisionflex.Consideration{
			slideIs(DRY),
			noDryScheduled,
			tipIs(FULL),
			robotAt(DISPENSE),
		},
		Actions: []decisionflex.Action{catMutator{func(c *Context) {
			c.SlideIs = WET
			c.TipIs = DIRTY
		}}},
	}

	goToShucker := decisionflex.Rule{
		Name: "go to tip shucker",
		Considerations: []decisionflex.Consideration{
			tipIs(DIRTY),
			robotNotAt(SHUCKER),
		},
		Actions: []decisionflex.Action{goTo(SHUCKER)},
	}

	shuckTip := decisionflex.Rule{
		Name: "shuck the tip",
		Considerations: []decisionflex.Consideration{
			tipIs(DIRTY),
			robotAt(SHUCKER),
		},
		Actions: []decisionflex.Action{catMutator{func(c *Context) {
			c.TipIs = NO_TIP
		}}},
	}

	acquireWet := decisionflex.Rule{
		Name: "acquire a wet reading",
		Considerations: []decisionflex.Consideration{
			slideIs(WET),
			someWetScheduled,
		},
		Actions: []decisionflex.Action{catMutator{func(c *Context) {
			c.WetScheduled--
		}}},
	}

	goToEject := decisionflex.Rule{
		Name: "go to slide eject station",
		Considerations: []decisionflex.Consideration{
			slideIs(WET),
			tipIs(NO_TIP),
			robotNotAt(EJECT),
		},
		Actions: []decisionflex.Action{goTo(EJECT)},
	}

	eject := decisionflex.Rule{
		Name: "eject a slide",
		Considerations: []decisionflex.Consideration{
			noWetScheduled,
			tipIs(NO_TIP),
			robotAt(EJECT),
		},
		Actions: []decisionflex.Action{catMutator{func(c *Context) {
			c.SlideIs = NO_SLIDE
		}}},
	}

	acquireDry := decisionflex.Rule{
		Name: "acquire a dry reading",
		Considerations: []decisionflex.Consideration{
			slideIs(DRY),
			someDryScheduled,
		},
		Actions: []decisionflex.Action{catMutator{func(c *Context) {
			c.DryScheduled--
		}}},
	}

	idle := decisionflex.Rule{
		Name: "idle",
		Considerations: []decisionflex.Consideration{
			decisionflex.ScalarConsideration{0.0},
		},
		Actions: []decisionflex.Action{},
	}

	seedPtr := flag.Int64(
		"seed",
		time.Now().UnixNano(), // default value
		"seed for random number generator",
	)
	flag.Parse()

	// randPtr := rand.New(rand.NewSource(*seedPtr))
	_ = rand.New(rand.NewSource(*seedPtr))

	deciderOut := log.New(os.Stdout, "", 0)
	decider := decisionflex.DecisionFlex{
		ContextFactory: decisionflex.SingleContextFactory{&context},
		// Chooser: decisionflex.FirstPossible,
		// Chooser: decisionflex.UniformRandom{*randPtr},
		// Chooser: decisionflex.WeightedRandom{*randPtr},
		// Chooser: decisionflex.SoftMax{0.5, *randPtr},
		Chooser: decisionflex.OptimalStopping{1},
		Rules: []decisionflex.Rule{
			acquireDry,
			acquireWet,
			goToTip,
			loadTip,
			goToSample,
			aspirate,
			goToDispense,
			dispense,
			goToShucker,
			shuckTip,
			goToEject,
			eject,
			idle,
		},
		Logger: *deciderOut,
	}

	for i := 0; context.SlideIs != NO_SLIDE && i < 100; i++ {
		decider.PerformAction()
	}
}
