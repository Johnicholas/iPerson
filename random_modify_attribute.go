// RandomModifyAttribute is a class, which is an Action
//
// It implements Perform (a method that takes an IContext)
// by pulling a random float from GetNextModify (a uniform random distribution between min and max)
// and invoking the target's BoostAttribute method with it

package main
