package mesos

import (
	"log"
)

// We are reached here only from the 'Ready' state
func stateHeartbeat(m *MesosMaster) stateFn {
	log.Print("STATE: Heartbeat")

	return stateReady
}
