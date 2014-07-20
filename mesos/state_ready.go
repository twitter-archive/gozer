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

	case command := <-m.command:
		stateSendCommand := func(fm *MesosMaster) stateFn {
			if err := command(fm); err != nil {
				log.Print("Error: ", err)
				return stateError
			}
			return stateReady
		}
		return stateSendCommand

	case event := <-m.events:
		stateReceiveEvent := func(fm *MesosMaster) stateFn {
			if err := fm.eventDispatch(event); err != nil {
				log.Print("Error: ", err)
				return stateError
			}
			return stateReady
		}
		return stateReceiveEvent
	}
}
