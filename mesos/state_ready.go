package mesos

import (
	"time"
)

const heartbeatTime = time.Minute

func stateReady(d *Driver) stateFn {
	// Framework is connected, ready and waiting for something to do
	d.config.Log.Debug.Println("STATE: Ready")

	select {
	case <-time.Tick(heartbeatTime):
		return stateHeartbeat

	case command, ok := <-d.command:
		if !ok {
			// stop when the command channel is closed
			return stateStop
		}
		stateSendCommand := func(fm *Driver) stateFn {
			if err := command(fm); err != nil {
				d.config.Log.Error.Println("Failed to run command:", err)
				return stateError
			}
			return stateReady
		}
		return stateSendCommand

	case event, ok := <-d.events:
		if !ok {
			// stop when the events channel is closed
			return stateStop
		}
		stateReceiveEvent := func(fm *Driver) stateFn {
			if err := fm.eventDispatch(event); err != nil {
				d.config.Log.Error.Println("Failed to dispatch event:", err)
				return stateError
			}
			return stateReady
		}
		return stateReceiveEvent
	}
}
