package main

import "fmt"
import "github.com/johnicholas/decisionflex"

func main() {
	dryScheduled := personAttribute{0.03}
	wetScheduled := personAttribute{0.21}
	atEject := personAttribute{0.0}
	atSlide := personAttribute{1.0}
	dry := personAttribute{1.0}
	wet := personAttribute{0.0}
	clean := personAttribute{0.0}
	full := personAttribute{0.0}
	dirty := personAttribute{0.0}
	atTip := personAttribute{0.0}
	atSample := personAttribute{0.0}
	atDispense := personAttribute{0.0}
	atShucker := personAttribute{0.0}

	/*
		loadSlide := []decisionflex.Performer{
			modifyAttribute{&dryScheduled, 0.03},
			modifyAttribute{&wetScheduled, 0.21},
			modifyAttribute{&dry, 1.0},
		}
	*/

	acquireDry := modifyAttribute{&dryScheduled, -0.01}
	acquireWet := modifyAttribute{&wetScheduled, -0.01}
	eject := modifyAttribute{&wet, -1.0}

	dispenseOn := []decisionflex.Performer{
		modifyAttribute{&dry, -1.0},
		modifyAttribute{&wet, 1.0},
		modifyAttribute{&full, -1.0},
		modifyAttribute{&dirty, 1.0},
	}
	loadTip := modifyAttribute{&clean, 1.0}
	aspirate := []decisionflex.Performer{
		modifyAttribute{&clean, -1.0},
		modifyAttribute{&full, 1.0},
	}
	shuckTip := modifyAttribute{&dirty, -1.0}
	leaveSlide := modifyAttribute{&atSlide, -1.0}
	arriveTip := modifyAttribute{&atTip, 1.0}
	leaveTip := modifyAttribute{&atTip, -1.0}
	arriveSlide := modifyAttribute{&atSlide, 1.0}
	arriveSample := modifyAttribute{&atSample, 1.0}
	leaveSample := modifyAttribute{&atSample, -1.0}
	arriveDispense := modifyAttribute{&atDispense, 1.0}
	leaveDispense := modifyAttribute{&atDispense, -1.0}
	arriveShucker := modifyAttribute{&atShucker, 1.0}
	leaveShucker := modifyAttribute{&atShucker, -1.0}
	arriveEject := modifyAttribute{&atEject, 1.0}

	someDryScheduled := decisionflex.BinaryContextConsiderer{
		ContextName:    "dry_scheduled",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	noDryScheduled := decisionflex.BinaryContextConsiderer{
		ContextName:    "dry_scheduled",
		EvaluationType: decisionflex.TrueIfEqualTo,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	isDry := decisionflex.BinaryContextConsiderer{
		ContextName:    "dry",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	robotAtSlide := decisionflex.BinaryContextConsiderer{
		ContextName:    "at_slide",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	notClean := decisionflex.BinaryContextConsiderer{
		ContextName:    "clean",
		EvaluationType: decisionflex.TrueIfEqualTo,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	notFull := decisionflex.BinaryContextConsiderer{
		ContextName:    "full",
		EvaluationType: decisionflex.TrueIfEqualTo,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	robotAtTip := decisionflex.BinaryContextConsiderer{
		ContextName:    "at_tip",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	isClean := decisionflex.BinaryContextConsiderer{
		ContextName:    "clean",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	robotAtSample := decisionflex.BinaryContextConsiderer{
		ContextName:    "at_sample",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	isFull := decisionflex.BinaryContextConsiderer{
		ContextName:    "full",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	robotAtDispense := decisionflex.BinaryContextConsiderer{
		ContextName:    "at_dispense",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	isDirty := decisionflex.BinaryContextConsiderer{
		ContextName:    "dirty",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	robotAtShucker := decisionflex.BinaryContextConsiderer{
		ContextName:    "at_shucker",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	isWet := decisionflex.BinaryContextConsiderer{
		ContextName:    "wet",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	someWetScheduled := decisionflex.BinaryContextConsiderer{
		ContextName:    "wet_scheduled",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	robotAtEject := decisionflex.BinaryContextConsiderer{
		ContextName:    "at_eject",
		EvaluationType: decisionflex.TrueIfMoreThan,
		GateValue:      0.0,
		IfTrue:         1.0,
		IfFalse:        0.0,
	}
	/*
		noWetScheduled := decisionflex.BinaryContextConsiderer{
			ContextName:    "wet_scheduled",
			EvaluationType: decisionflex.TrueIfEqualTo,
			GateValue:      0.0,
			IfTrue:         1.0,
			IfFalse:        0.0,
		}
	*/

	// if (dry>0&&dry_scheduled>0) fire(acquire_dry)
	possiblyAcquireDry := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly acquire dry"},
		Considerations: []decisionflex.Considerer{isDry, someDryScheduled},
		Actions:        []decisionflex.Performer{acquireDry},
	}
	// if (dry>0&&dry_scheduled==0&&at_slide>0&&clean==0) fire(at_tip)
	possiblyGoToTip := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly go to tip"},
		Considerations: []decisionflex.Considerer{isDry, noDryScheduled, robotAtSlide, notClean},
		Actions:        []decisionflex.Performer{leaveSlide, arriveTip},
	}
	// if (dry>0&&dry_scheduled==0&&clean==0&&full==0&&at_tip>0) fire(load_tip)
	possiblyLoadTip := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly load tip"},
		Considerations: []decisionflex.Considerer{isDry, noDryScheduled, notClean, notFull, robotAtTip},
		Actions:        []decisionflex.Performer{loadTip},
	}
	// if(clean>0&&full==0&&at_tip>0) fire(at_slide)
	possiblyGoToSlide := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly go to slide"},
		Considerations: []decisionflex.Considerer{isClean, notFull, robotAtTip},
		Actions:        []decisionflex.Performer{leaveTip, arriveSlide},
	}
	// if (clean>0&&full==0&&at_slide>0) fire(at_sample)
	possiblyGoToSample := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly go to sample"},
		Considerations: []decisionflex.Considerer{isClean, notFull, robotAtSlide},
		Actions:        []decisionflex.Performer{leaveSlide, arriveSample},
	}
	// if (clean>0&&full==0&&at_sample>0) fire(aspirate)
	possiblyAspirate := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly aspirate"},
		Considerations: []decisionflex.Considerer{isClean, notFull, robotAtSample},
		Actions:        aspirate,
	}
	// if (dry>0&&full>0&&at_sample>0) fire(at_dispense)
	possiblyGoToDispense := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly go to dispense"},
		Considerations: []decisionflex.Considerer{isDry, isFull, robotAtSample},
		Actions:        []decisionflex.Performer{leaveSample, arriveDispense},
	}
	// if (dry>0&&full>0&&at_dispense) fire(dispense_on)
	possiblyDispense := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly dispense sample onto slide"},
		Considerations: []decisionflex.Considerer{isDry, isFull, robotAtDispense},
		Actions:        dispenseOn,
	}
	// if (dirty>0&&at_dispense>0) fire(at_shucker)
	possiblyGoToShucker := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly go to shucker"},
		Considerations: []decisionflex.Considerer{isDirty, robotAtDispense},
		Actions:        []decisionflex.Performer{leaveDispense, arriveShucker},
	}
	// if (dirty>0&&at_shucker>0) fire(shuck_tip)
	possiblyShuckTip := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly shuck the tip"},
		Considerations: []decisionflex.Considerer{isDirty, robotAtShucker},
		Actions:        []decisionflex.Performer{shuckTip},
	}
	// if (wet>0&&wet_schedule>0) fire(acquire_wet)
	possiblyAcquireWet := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly acquire wet"},
		Considerations: []decisionflex.Considerer{isWet, someWetScheduled},
		Actions:        []decisionflex.Performer{acquireWet},
	}
	// if (wet>0&&at_shucker>0) fire(at_eject)
	possiblyGoToEject := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly go to eject"},
		Considerations: []decisionflex.Considerer{isWet, robotAtShucker},
		Actions:        []decisionflex.Performer{leaveShucker, arriveEject},
	}
	// if (wet>0&&at_eject>0) fire(eject)
	possiblyEject := decisionflex.ActionConsiderations{
		ActionsObject:  decisionflex.GameObject{"possibly eject a slide"},
		Considerations: []decisionflex.Considerer{isWet, robotAtEject},
		Actions:        []decisionflex.Performer{eject},
	}

	decider := decisionflex.Decisionflex{
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
		Enabled: true,
	}
	answer := decider.PerformAction()

	fmt.Println(answer.Score)
	fmt.Println(answer.ActionObject.Name)
}
