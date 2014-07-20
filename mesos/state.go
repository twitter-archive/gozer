package mesos

import (
	"log"
)

// A state function is a function that returns a state function
type stateFn func(*MesosMaster) stateFn

// Run the state machine
func (m *MesosMaster) Run() {
	for state := stateInit; state != nil; {
		state = state(m)
	}
}

func stateInit(m *MesosMaster) stateFn {
	log.Print("INIT: Starting framework:", m)

	return stateRegistering
}

func stateStop(m *MesosMaster) stateFn {
	log.Print("STOP: Stopping framework:", m)
	return nil
}

func stateRegistering(m *MesosMaster) stateFn {
	log.Print("REGISTERING: Trying to register framework:", m)

	return stateReady
}
