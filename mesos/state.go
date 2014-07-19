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
	// ...
	return stateStop
}

func stateStop(m *MesosMaster) stateFn {
	log.Print("STOP: Stopping framework:", m)
	return nil
}

func stateRegistering(m *MesosMaster) stateFn {
	// ...
	return nil
}

func stateSendCommand(m *MesosMaster) stateFn {
	// ...
	return nil
}

func stateReceiveEvent(m *MesosMaster) stateFn {
	// ...
	return nil
}
