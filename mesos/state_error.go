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

func stateError(d *Driver) stateFn {
	// Handle any type of error state
	log.Print("STATE: Error, MesosMaster = ", d)
	return stateStop
}
