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
	EJECT    Location = iota
	SLIDE             = iota
	TIP               = iota
	SAMPLE            = iota
	DISPENSE          = iota
	SHUCKER           = iota
)

type SlideState int

const (
	NO_SLIDE SlideState = iota
	DRY                 = iota
	WET                 = iota
)

type TipState int

const (
	NO_TIP TipState = iota
	CLEAN           = iota
	FULL            = iota
	DIRTY           = iota
)

type Context struct {
	DryScheduled int
	WetScheduled int
	RobotAt      Location
	SlideIs      SlideState
	TipIs        TipState
}

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

func someDryScheduled(c Context) bool {
	return c.DryScheduled > 0
}

func noDryScheduled(c Context) bool {
	return c.DryScheduled == 0
}

func someWetScheduled(c Context) bool {
	return c.WetScheduled > 0
}

func noWetScheduled(c Context) bool {
	return c.WetScheduled == 0
}

func addPrecondition(r *decisionflex.Rule, toAdd func(Context) bool) {
	r.Considerations = append(r.Considerations, catPrecondition{toAdd})
}

type catMutator struct {
	toRun func(*Context)
}

func (c catMutator) Perform(i interface{}) {
	downcast := i.(*Context)
	c.toRun(downcast)
}

func goTo(l Location) func(c *Context) {
	return func(c *Context) {
		c.RobotAt = l
	}
}

func robotAt(l Location) func(Context) bool {
	return func(c Context) bool {
		return c.RobotAt == l
	}
}

func robotNotAt(l Location) func(Context) bool {
	return func(c Context) bool {
		return c.RobotAt != l
	}
}

func slideIs(s SlideState) func(Context) bool {
	return func(c Context) bool {
		return c.SlideIs == s
	}
}

func tipIs(t TipState) func(Context) bool {
	return func(c Context) bool {
		return c.TipIs == t
	}
}

func tipIsNot(t TipState) func(Context) bool {
	return func(c Context) bool {
		return c.TipIs != t
	}
}

func addEffect(r *decisionflex.Rule, m func(*Context)) {
	r.Actions = append(r.Actions, catMutator{m})
}

func addRule(d *decisionflex.DecisionFlex, r decisionflex.Rule) {
	d.Rules = append(d.Rules, r)
}

func main() {
	context := Context{
		DryScheduled: 3,
		WetScheduled: 21,
		RobotAt:      SLIDE,
		SlideIs:      DRY,
		// TipIs: NO_TIP
	}

	// if (dry>0&&dry_scheduled==0&&at_slide>0&&clean==0) fire(at_tip)
	goToTip := decisionflex.Rule{Name: "go to tip load station"}
	addPrecondition(&goToTip, slideIs(DRY))
	addPrecondition(&goToTip, tipIs(NO_TIP))
	addPrecondition(&goToTip, robotNotAt(TIP))
	addEffect(&goToTip, goTo(TIP))

	// if (dry>0&&dry_scheduled==0&&tip_none>0&&at_tip>0) fire(load_tip)
	loadTip := decisionflex.Rule{Name: "load a tip"}
	addPrecondition(&loadTip, slideIs(DRY))
	addPrecondition(&loadTip, tipIs(NO_TIP))
	addPrecondition(&loadTip, robotAt(TIP))
	addEffect(&loadTip, func(c *Context) {
		// preconditions
		if c.RobotAt != TIP {
			panic("load tip when robot is not at tip load station")
		}
		if c.TipIs != NO_TIP {
			panic("load tip with another tip already on proboscis")
		}
		// effects
		c.TipIs = CLEAN
	})

	// if (clean>0&&at_sample==0) fire(at_sample)
	goToSample := decisionflex.Rule{Name: "go to sample cup"}
	addPrecondition(&goToSample, slideIs(DRY))
	addPrecondition(&goToSample, tipIs(CLEAN))
	addPrecondition(&goToSample, robotNotAt(SAMPLE))
	addEffect(&goToSample, goTo(SAMPLE))

	// if (clean>0&&at_sample>0) fire(aspirate)
	aspirate := decisionflex.Rule{Name: "aspirate from sample cup"}
	addPrecondition(&aspirate, slideIs(DRY))
	addPrecondition(&aspirate, noDryScheduled)
	addPrecondition(&aspirate, tipIs(CLEAN))
	addPrecondition(&aspirate, robotAt(SAMPLE))
	addEffect(&aspirate, func(c *Context) {
		// preconditions
		if c.RobotAt != SAMPLE {
			panic("aspirate while robot is not at sample cup")
		}
		if c.TipIs != CLEAN {
			panic("aspirate while tip is not clean")
		}
		// effects
		c.TipIs = FULL
	})

	// if (dry>0&&full>0&&at_dispense==0) fire(at_dispense)
	goToDispense := decisionflex.Rule{Name: "go to dispense station"}
	addPrecondition(&goToDispense, slideIs(DRY))
	addPrecondition(&goToDispense, noDryScheduled)
	addPrecondition(&goToDispense, tipIs(FULL))
	addPrecondition(&goToDispense, robotNotAt(DISPENSE))
	addEffect(&goToDispense, goTo(DISPENSE))

	// if (dry>0&&full>0&&at_dispense) fire(dispense_on)
	dispense := decisionflex.Rule{Name: "dispense sample onto slide"}
	addPrecondition(&dispense, slideIs(DRY))
	addPrecondition(&dispense, noDryScheduled)
	addPrecondition(&dispense, tipIs(FULL))
	addPrecondition(&dispense, robotAt(DISPENSE))
	addEffect(&dispense, func(c *Context) {
		// preconditions
		if c.SlideIs != DRY {
			panic("dispense on slide when slide is not dry")
		}
		if c.TipIs != FULL {
			panic("dispense on slide when tip is not full")
		}
		// effects
		c.SlideIs = WET
		c.TipIs = DIRTY
	})

	// if (dirty>0&&at_dispense>0) fire(at_shucker)
	goToShucker := decisionflex.Rule{Name: "go to tip shucker"}
	addPrecondition(&goToShucker, tipIs(DIRTY))
	addPrecondition(&goToShucker, robotNotAt(SHUCKER))
	addEffect(&goToShucker, goTo(SHUCKER))

	// if (dirty>0&&at_shucker>0) fire(shuck_tip)
	shuckTip := decisionflex.Rule{Name: "shuck the tip"}
	addPrecondition(&shuckTip, tipIs(DIRTY))
	addPrecondition(&shuckTip, robotAt(SHUCKER))
	addEffect(&shuckTip, func(c *Context) {
		// preconditions
		if c.RobotAt != SHUCKER {
			panic("shuck tip while robot is not at tip shucker")
		}
		// effects
		c.TipIs = NO_TIP
	})

	// if (wet>0&&wet_schedule>0) fire(acquire_wet)
	acquireWet := decisionflex.Rule{Name: "acquire a wet reading"}
	addPrecondition(&acquireWet, slideIs(WET))
	addPrecondition(&acquireWet, someWetScheduled)
	addEffect(&acquireWet, func(c *Context) {
		// preconditions
		if c.WetScheduled <= 0 {
			panic("unscheduled acquire wet")
		}
		// effects
		c.WetScheduled--
	})

	// if (wet>0&&no_tip>0&&at_eject==0) fire(at_eject)
	goToEject := decisionflex.Rule{Name: "go to eject station"}
	addPrecondition(&goToEject, slideIs(WET))
	addPrecondition(&goToEject, tipIs(NO_TIP))
	addPrecondition(&goToEject, robotNotAt(EJECT))
	addEffect(&goToEject, goTo(EJECT))

	// if (wet_scheduled==0&&no_tip>0&&at_eject>0) fire(eject)
	eject := decisionflex.Rule{Name: "eject a slide"}
	addPrecondition(&eject, noWetScheduled)
	addPrecondition(&eject, tipIs(NO_TIP))
	addPrecondition(&eject, robotAt(EJECT))
	addEffect(&eject, func(c *Context) {
		// preconditions
		if c.RobotAt != EJECT {
			panic("eject when robot is not at eject station")
		}
		if c.TipIs != NO_TIP {
			panic("eject when proboscis has a tip on it")
		}
		// effect
		c.SlideIs = NO_SLIDE
	})

	// if (dry>0&&dry_scheduled>0) fire(acquire_dry)
	acquireDry := decisionflex.Rule{Name: "acquire a dry reading"}
	addPrecondition(&acquireDry, slideIs(DRY))
	addPrecondition(&acquireDry, someDryScheduled)
	addEffect(&acquireDry, func(c *Context) {
		// preconditions
		if c.SlideIs != DRY {
			panic("cannot take a dry read when slide is not dry")
		}
		if c.DryScheduled <= 0 {
			panic("unscheduled acquire dry")
		}
		// effects
		c.DryScheduled--
	})

	idle := decisionflex.Rule{
		Name: "idle",
		Considerations: []decisionflex.Consideration{
			decisionflex.ScalarConsideration{0.0},
		},
		Actions: []decisionflex.Performer{},
	}

	seedPtr := flag.Int64(
		"seed",
		time.Now().UnixNano(), // default value
		"seed for random number generator",
	)
	flag.Parse()

	randPtr := rand.New(rand.NewSource(*seedPtr))

	deciderOut := log.New(os.Stdout, "", 0)
	decider := decisionflex.DecisionFlex{
		ContextFactory: decisionflex.SingleContextFactory{&context},
		// Chooser:        decisionflex.FirstPossible,
		// Chooser: decisionflex.UniformRandom{*randPtr},
		// Chooser: decisionflex.WeightedRandom{*randPtr},
		Chooser: decisionflex.SoftMax{1.5, *randPtr},
		Logger:  *deciderOut,
	}

	addRule(&decider, acquireDry)
	addRule(&decider, acquireWet)

	addRule(&decider, goToTip)
	addRule(&decider, loadTip)

	addRule(&decider, goToSample)
	addRule(&decider, aspirate)

	addRule(&decider, goToDispense)
	addRule(&decider, dispense)

	addRule(&decider, goToShucker)
	addRule(&decider, shuckTip)

	addRule(&decider, goToEject)
	addRule(&decider, eject)

	addRule(&decider, idle)

	for i := 0; context.SlideIs != NO_SLIDE && i < 100; i++ {
		decider.PerformAction()
	}
}
