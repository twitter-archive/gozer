package mesos

import (
	"log"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

func stateRegister(m *MesosMaster) stateFn {
	log.Print("REGISTERING: Trying to register framework:", m)

	// Create the register message and send it.
	callType := mesos_scheduler.Call_REGISTER
	registerCall := &mesos_scheduler.Call{
		Type: &callType,
		FrameworkInfo: &mesos.FrameworkInfo{
			User: &m.config.RegisteredUser,
			Name: &m.config.FrameworkName,
		},
	}

	// TODO(weingart): This should re-try and backoff
	err := m.send(registerCall)
	if err != nil {
		log.Print("Error: send: ", err)
		return stateError
	}

	// Wait for Registered event, throw away any other events
	for {
		event := <-m.events
		if *event.Type != mesos_scheduler.Event_REGISTERED {
			log.Printf("Unexpected event type: want %q, got %+v", mesos_scheduler.Event_REGISTERED, *event.Type)
		}
		m.frameworkId = *event.Registered.FrameworkId
		break
	}

	log.Printf("Registered %s:%s with id %q", m.config.RegisteredUser, m.config.FrameworkName, *m.frameworkId.Value)
	return stateReady
}
