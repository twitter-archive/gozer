package mesos

import (
	"time"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

const maxRegisterWait = 20 * time.Second

func stateRegister(d *Driver) stateFn {
	d.config.Log.Info.Printf("REGISTERING: Trying to register framework: %+v", d)

	// Create the register message and send it.
	callType := mesos_scheduler.Call_REGISTER
	registerCall := &mesos_scheduler.Call{
		Type: &callType,
		FrameworkInfo: &mesos.FrameworkInfo{
			User: &d.config.RegisteredUser,
			Name: &d.config.FrameworkName,
		},
	}

	err := d.send(registerCall)
	registerBackoff := 1 * time.Second

	if err != nil {
		d.config.Log.Warn.Println("Failed to send register:", err)
		d.config.Log.Warn.Printf("Waiting for %s before trying again.", registerBackoff)
	}

	// Wait for Registered event, throw away any other events
	registered := false
	for !registered {
		select{
		case event := <-d.events:
			if *event.Type != mesos_scheduler.Event_REGISTERED {
				d.config.Log.Error.Printf("Unexpected event type: want %q, got %+v",
					mesos_scheduler.Event_REGISTERED, *event.Type)
				continue
			}
			d.frameworkId = *event.Registered.FrameworkId
			registered = true
			break

		case <-time.After(maxRegisterWait):
			d.config.Log.Error.Printf("Failed to register after %s", maxRegisterWait)
			return stateError

		case <-time.After(registerBackoff):
			err := d.send(registerCall)
			if err != nil {
				registerBackoff = registerBackoff * 2
				d.config.Log.Warn.Println("Failed to send register:", err)
				d.config.Log.Warn.Printf("Waiting for %s before trying again.", registerBackoff)
			}
		}
	}

	d.config.Log.Info.Printf("Registered %s:%s with id %q",
		d.config.RegisteredUser, d.config.FrameworkName, *d.frameworkId.Value)
	return stateReady
}
