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

	case cmd := <-m.command:
		stateSendCommand := func(fm *MesosMaster) stateFn {
			if err := cmd(fm); err != nil {
				log.Print("Error:", err)
				return stateError
			}
			return stateReady
		}
		return stateSendCommand

	case <-m.events:
		stateReceiveEvent := func(fm *MesosMaster) stateFn {

			return stateReady
		}
		return stateReceiveEvent
	}
}
