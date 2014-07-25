package mesos

import (
	"log"
	"time"
)

func stateReady(d *Driver) stateFn {
	// Framework is connected, ready and waiting for something to do
	log.Print("STATE: Ready")

	select {
	case <-time.Tick(time.Minute):
		return stateHeartbeat

	case command := <-d.command:
		stateSendCommand := func(fm *Driver) stateFn {
			if err := command(fm); err != nil {
				log.Print("Error: ", err)
				return stateError
			}
			return stateReady
		}
		return stateSendCommand

	case event := <-d.events:
		stateReceiveEvent := func(fm *Driver) stateFn {
			if err := fm.eventDispatch(event); err != nil {
				log.Print("Error: ", err)
				return stateError
			}
			return stateReady
		}
		return stateReceiveEvent
	}
}
