package mesos

import (
	"log"
	"time"
)

const heartbeatTime = time.Minute

func stateReady(d *Driver) stateFn {
	// Framework is connected, ready and waiting for something to do
	log.Println("STATE: Ready")

	select {
	case <-time.Tick(heartbeatTime):
		return stateHeartbeat

	case command := <-d.command:
		stateSendCommand := func(fm *Driver) stateFn {
			if err := command(fm); err != nil {
				log.Print("Error running command: ", err)
				return stateError
			}
			return stateReady
		}
		return stateSendCommand

	case event := <-d.events:
		stateReceiveEvent := func(fm *Driver) stateFn {
			if err := fm.eventDispatch(event); err != nil {
				log.Print("Error dispatching event: ", err)
				return stateError
			}
			return stateReady
		}
		return stateReceiveEvent
	}
}
