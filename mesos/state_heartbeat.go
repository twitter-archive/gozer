package mesos

import (
	"log"
)

// We are reached here only from the 'Ready' state
func stateHeartbeat(d *Driver) stateFn {
	log.Println("STATE: Heartbeat")
	return stateReady
}
