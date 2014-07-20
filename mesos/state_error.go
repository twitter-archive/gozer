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
	// Handle any type of error state
	log.Print("STATE: Error, MesosMaster = ", m)
	return stateStop
}
