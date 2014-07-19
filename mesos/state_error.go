package mesos

import (
	"log"
)

type mesosError int

const (
	errorNone mesosError = iota
	errorNotInitialized
	errorNotConnected
	errorNotReady
	//...
)

func stateError(m *MesosMaster) stateFn {
	m.RLock()
	defer m.Unlock()

	// Handle any type of error state
	log.Print("STATE: Error, MesosMaster = ", m)
	return stateStop
}
