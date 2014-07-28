package mesos

import (
	"log"
)

// A state function is a function that does stuff, and
// then returns the next state function to be invoked.
type stateFn func(*Driver) stateFn

// Run the state machine
func (d *Driver) Run() {
	for state := stateInit; state != nil; {
		state = state(d)
	}
}

func stateStop(d *Driver) stateFn {
	log.Print("STOP: Stopping framework: ", d)
	return nil
}