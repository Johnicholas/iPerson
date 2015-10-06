package main

import "fmt"
import "github.com/johnicholas/decisionflex"

func main() {
	dryScheduled := new(personAttribute)
	wetScheduled := new(personAttribute)
	dry := new(personAttribute)
	wet := new(personAttribute)
	clean := new(personAttribute)
	full := new(personAttribute)
	dirty := new(personAttribute)

	loadSlide := [...]decisionflex.Performer{
		modifyAttribute{dryScheduled, 0.03},
		modifyAttribute{wetScheduled, 0.21},
		modifyAttribute{dry, 1.0},
	}
	acquireDry := modifyAttribute{dryScheduled, -1.0}
	acquireWet := modifyAttribute{wetScheduled, -1.0}
	eject := modifyAttribute{wet, -1.0}
	dispenseOn := [...]decisionflex.Performer{
		modifyAttribute{dry, -1.0},
		modifyAttribute{wet, 1.0},
		modifyAttribute{full, -1.0},
		modifyAttribute{dirty, 1.0},
	}
	loadTip := modifyAttribute{clean, 1.0}
	aspirate := [...]decisionflex.Performer{
		modifyAttribute{clean, -1.0},
		modifyAttribute{full, 1.0},
	}
	shuckTip := modifyAttribute{dirty, -1.0}

        possibly_acquire_dry := decisionflex.ActionConsiderations{
          [...]decisionflex.Performer{modifyAttribute{dry_scheduled, -1.0}},
          [...]decisionflex.Considerer{dry_scheduled, dry},
          gameObject{"possibly_acquire_dry"},
        }
        possibly_at_tip := decisionflex.ActionConsiderations{
          [...]decisionflex.Performer{modifyAttribute{at_tip, 1.0}},
          [...]decisionflex.Considerer{
            dry,
            dry_scheduled == 0,
            at_slide>0,
            clean==0,
          },
          gameObject{"possibly_at_tip"},
        }

/*
if (dry>0&&dry_scheduled>0) fire(acquire_dry)
if (dry>0&&dry_scheduled==0&&at_slide>0&&clean==0) fire(at_tip)
if (dry>0&&dry_scheduled==0&&clean==0&&full==0&&at_tip>0) fire(load_tip)
if(clean>0&&full==0&&at_tip>0)fire(at_slide)
if (clean>0&&full==0&&at_slide>0)fire(at_sample)
if (clean>0&&full==0&&at_sample>0) fire(aspirate)
if (dry>0&&full>0&&at_sample>0) fire(at_dispense)
if (dry>0&&full>0&&at_dispense) fire(dispense_on)
if (dirty>0&&at_dispense>0) fire(at_shucker)
if (dirty>0&&at_shucker>0) fire(shuck_tip)
if (wet>0&&wet_scheduled>0) fire(acquire_wet)
if (wet>0&&at_shucker>0) fire(at_eject)
if (wet>0&&at_eject>0) fire(eject)
*/

	fmt.Printf("Hello, world.\n")
}
