package mesos

import (
	"log"
	"time"
)

func stateReady(m *MesosMaster) stateFn {
	// Framework is connected, ready and waiting for something to do
	log.Print("STATE: Ready")

	select {
	case <-time.Tick(time.Minute):
		return stateHeartbeat

	case <-m.sendCommand:
		return stateSendCommand

	case <-m.receivedEvent:
		return stateReceiveEvent
	}
}
