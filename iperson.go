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

type acquireDry struct{}

func (a acquireDry) Perform(context interface{}) {
	context.(*catContext).DryScheduled--
}

type acquireWet struct{}

func (a acquireWet) Perform(context interface{}) {
	context.(*catContext).WetScheduled--
}

type eject struct{}

func (e eject) Perform(context interface{}) {
	context.(*catContext).SlideIs = SLIDE_NONE
}

type dispenseOn struct{}

func (d dispenseOn) Perform(context interface{}) {
	// expect DRY
	context.(*catContext).SlideIs = WET
	// expect CLEAN
	context.(*catContext).TipIs = DIRTY
}

type loadTip struct{}

func (l loadTip) Perform(context interface{}) {
	// expect RobotIs == TIP
	// expect TIP_NONE
	context.(*catContext).TipIs = CLEAN
}

type aspirate struct{}

func (a aspirate) Perform(context interface{}) {
	// expect RobotIs == SAMPLE
	// expect TipIs == CLEAN
	context.(*catContext).TipIs = FULL
}

type shuckTip struct{}

func (a shuckTip) Perform(context interface{}) {
	context.(*catContext).TipIs = TIP_NONE
}

type goTo struct {
	Destination int
}

func (g goTo) Perform(context interface{}) {
	context.(*catContext).RobotAt = g.Destination
}

type someDryScheduled struct{}

func (s someDryScheduled) Consider(context interface{}) float64 {
	if context.(*catContext).DryScheduled > 0 {
		return 1.0
	} else {
		return 0.0
	}
}

type noDryScheduled struct{}

func (s noDryScheduled) Consider(context interface{}) float64 {
	if context.(*catContext).DryScheduled == 0 {
		return 1.0
	} else {
		return 0.0
	}
}

type someWetScheduled struct{}

func (s someWetScheduled) Consider(context interface{}) float64 {
	if context.(*catContext).WetScheduled > 0 {
		return 1.0
	} else {
		return 0.0
	}
}

type noWetScheduled struct{}

func (s noWetScheduled) Consider(context interface{}) float64 {
	if context.(*catContext).WetScheduled == 0 {
		return 1.0
	} else {
		return 0.0
	}
}

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

func main() {
	context := catContext{
		DryScheduled: 3,
		WetScheduled: 21,
		RobotAt:      SLIDE,
		// SlideIs: NONE,
		// TipIs: NONE
	}

	// if (dry>0&&dry_scheduled>0) fire(acquire_dry)
	possiblyAcquireDry := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly acquire dry"},
		Considerations: []decisionflex.Considerer{slideIs{DRY}, someDryScheduled{}},
		Actions:        []decisionflex.Performer{acquireDry{}},
	}
	// if (dry>0&&dry_scheduled==0&&at_slide>0&&clean==0) fire(at_tip)
	possiblyGoToTip := decisionflex.ActionConsiderations{
		ActionsObject: decisionflex.GameObject{"possibly go to tip"},
		Considerations: []decisionflex.Considerer{
			slideIs{DRY},
			noDryScheduled{},
			robotAt{SLIDE},
			tipIsNot{CLEAN},
		},
		Actions: []decisionflex.Performer{goTo{TIP}},
	}
	// if (dry>0&&dry_scheduled==0&&clean==0&&full==0&&at_tip>0) fire(load_tip)
	possiblyLoadTip := decisionflex.ActionConsiderations{
		ActionsObject: decisionflex.GameObject{"possibly load tip"},
		Considerations: []decisionflex.Considerer{
			slideIs{DRY},
			noDryScheduled{},
			tipIsNot{CLEAN},
			tipIsNot{FULL},
			robotAt{TIP},
		},
		Actions: []decisionflex.Performer{loadTip{}},
	}
	// if(clean>0&&full==0&&at_tip>0) fire(at_slide)
	possiblyGoToSlide := decisionflex.ActionConsiderations{
		ActionsObject: decisionflex.GameObject{"possibly go to slide"},
		Considerations: []decisionflex.Considerer{
			tipIs{CLEAN},
			tipIsNot{FULL},
			robotAt{TIP},
		},
		Actions: []decisionflex.Performer{goTo{SLIDE}},
	}
	// if (clean>0&&full==0&&at_slide>0) fire(at_sample)
	possiblyGoToSample := decisionflex.ActionConsiderations{
		ActionsObject: decisionflex.GameObject{"possibly go to sample"},
		Considerations: []decisionflex.Considerer{
			tipIs{CLEAN},
			tipIsNot{FULL},
			robotAt{SLIDE},
		},
		Actions: []decisionflex.Performer{goTo{SAMPLE}},
	}
	// if (clean>0&&full==0&&at_sample>0) fire(aspirate)
	possiblyAspirate := decisionflex.ActionConsiderations{
		ActionsObject: decisionflex.GameObject{"possibly aspirate"},
		Considerations: []decisionflex.Considerer{
			tipIs{CLEAN},
			tipIsNot{FULL},
			robotAt{SAMPLE},
		},
		Actions: []decisionflex.Performer{aspirate{}},
	}
	// if (dry>0&&full>0&&at_sample>0) fire(at_dispense)
	possiblyGoToDispense := decisionflex.ActionConsiderations{
		ActionsObject: decisionflex.GameObject{"possibly go to dispense"},
		Considerations: []decisionflex.Considerer{
			slideIs{DRY},
			tipIs{FULL},
			robotAt{SAMPLE},
		},
		Actions: []decisionflex.Performer{goTo{DISPENSE}},
	}
	// if (dry>0&&full>0&&at_dispense) fire(dispense_on)
	possiblyDispense := decisionflex.ActionConsiderations{
		ActionsObject: decisionflex.GameObject{"possibly dispense sample onto slide"},
		Considerations: []decisionflex.Considerer{
			slideIs{DRY},
			tipIs{FULL},
			robotAt{DISPENSE},
		},
		Actions: []decisionflex.Performer{dispenseOn{}},
	}
	// if (dirty>0&&at_dispense>0) fire(at_shucker)
	possiblyGoToShucker := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly go to shucker"},
		Considerations: []decisionflex.Considerer{tipIs{DIRTY}, robotAt{DISPENSE}},
		Actions:        []decisionflex.Performer{goTo{SHUCKER}},
	}
	// if (dirty>0&&at_shucker>0) fire(shuck_tip)
	possiblyShuckTip := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly shuck the tip"},
		Considerations: []decisionflex.Considerer{tipIs{DIRTY}, robotAt{SHUCKER}},
		Actions:        []decisionflex.Performer{shuckTip{}},
	}
	// if (wet>0&&wet_schedule>0) fire(acquire_wet)
	possiblyAcquireWet := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly acquire wet"},
		Considerations: []decisionflex.Considerer{slideIs{WET}, someWetScheduled{}},
		Actions:        []decisionflex.Performer{acquireWet{}},
	}
	// if (wet>0&&at_shucker>0) fire(at_eject)
	possiblyGoToEject := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly go to eject"},
		Considerations: []decisionflex.Considerer{slideIs{WET}, robotAt{SHUCKER}},
		Actions:        []decisionflex.Performer{goTo{EJECT}},
	}
	// if (wet>0&&at_eject>0) fire(eject)
	possiblyEject := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly eject a slide"},
		Considerations: []decisionflex.Considerer{slideIs{WET}, robotAt{EJECT}},
		Actions:        []decisionflex.Performer{eject{}},
	}

	decider := decisionflex.DecisionFlex{
		Actions: []decisionflex.ActionConsiderations{
			possiblyAcquireDry,
			possiblyGoToTip,
			possiblyLoadTip,
			possiblyGoToSlide,
			possiblyGoToSample,
			possiblyAspirate,
			possiblyGoToDispense,
			possiblyDispense,
			possiblyGoToShucker,
			possiblyShuckTip,
			possiblyAcquireWet,
			possiblyGoToEject,
			possiblyEject,
		},
		Enabled:        true,
		Selector:       decisionflex.SelectWeightedRandom{0.0},
		ContextFactory: decisionflex.SingleContextFactory{&context},
	}

	answer := decider.PerformAction()

	fmt.Println(answer.Score)
	fmt.Println(answer.ActionObject.Name)
}
