package mesos

import (
	"log"
)

// A state function is a function that does stuff, and
// then returns the next state function to be invoked.
type stateFn func(*MesosMaster) stateFn

// Run the state machine
func (m *MesosMaster) Run() {
	for state := stateInit; state != nil; {
		state = state(m)
	}
}

func stateStop(m *MesosMaster) stateFn {
	log.Print("STOP: Stopping framework:", m)
	return nil
}
