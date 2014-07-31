package mesos

import (
	"time"
)

const heartbeatTime = time.Minute

func stateReady(d *Driver) stateFn {
	// Framework is connected, ready and waiting for something to do
	d.config.Log.Info.Println("STATE: Ready")

	select {
	case <-time.Tick(heartbeatTime):
		return stateHeartbeat

	case command := <-d.command:
		stateSendCommand := func(fm *Driver) stateFn {
			if err := command(fm); err != nil {
				d.config.Log.Error.Println("Error running command:", err)
				return stateError
			}
			return stateReady
		}
		return stateSendCommand

	case event := <-d.events:
		stateReceiveEvent := func(fm *Driver) stateFn {
			if err := fm.eventDispatch(event); err != nil {
				d.config.Log.Error.Println("Error dispatching event:", err)
				return stateError
			}
			return stateReady
		}
		return stateReceiveEvent
	}
}
