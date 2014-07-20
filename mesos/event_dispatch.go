package mesos

import (
	"fmt"
	"log"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

func (m *MesosMaster) eventDispatch(event *mesos_scheduler.Event) error {

	switch *event.Type {
	case mesos_scheduler.Event_REGISTERED:
		log.Print("Event REGISTERED: ", event)

	case mesos_scheduler.Event_REREGISTERED:
		log.Print("Event REREGISTERED: ", event)

	case mesos_scheduler.Event_OFFERS:
		for _, offer := range event.Offers.Offers {
			if *offer.FrameworkId.Value != *m.frameworkId.Value {
				log.Print("Unexpected framework in offer: want %q, got %q",
					*m.frameworkId.Value, *offer.FrameworkId.Value)
				continue
			}

			if len(m.Offers) < cap(m.Offers) {
				m.Offers <- offer
			} else {
				// TODO(weingart): how to ignore/return offer?
			}
		}

	case mesos_scheduler.Event_RESCIND:
		log.Print("Event RESCIND: ", event)

	case mesos_scheduler.Event_UPDATE:
		log.Print("Event UPDATE: ", event)

		acknowledgeType := mesos_scheduler.Call_ACKNOWLEDGE
		acknowledgeCall := &mesos_scheduler.Call{
			FrameworkInfo: &mesos.FrameworkInfo{
				User: &m.config.RegisteredUser,
				Name: &m.config.FrameworkName,
				Id:   &m.frameworkId,
			},
			Type: &acknowledgeType,
			Acknowledge: &mesos_scheduler.Call_Acknowledge{
				SlaveId: event.Update.Status.SlaveId,
				TaskId:  event.Update.Status.TaskId,
				Uuid:    event.Update.Uuid,
			},
		}

		err := m.send(acknowledgeCall)
		if err != nil {
			return fmt.Errorf("failed to send acknowledgement: %+v", err)
		}

		switch *event.Update.Status.State {
		case mesos.TaskState_TASK_STAGING:
		case mesos.TaskState_TASK_STARTING:
		case mesos.TaskState_TASK_RUNNING:
			log.Printf("task %s is running: %s", event.Update.Status.TaskId.GetValue(), event.Update.Status.GetMessage())
		case mesos.TaskState_TASK_FINISHED:
			log.Printf("task %s is complete: %s", event.Update.Status.TaskId.GetValue(), event.Update.Status.GetMessage())
		case mesos.TaskState_TASK_FAILED:
			log.Printf("task %s failed: %s", event.Update.Status.TaskId.GetValue(), event.Update.Status.GetMessage())
		case mesos.TaskState_TASK_KILLED:
		case mesos.TaskState_TASK_LOST:
			log.Printf("task %s failed to complete: %s", event.Update.Status.TaskId.GetValue(), event.Update.Status.GetMessage())
		}
		// TODO(weingart): Framework (gozer) should be informed here

	case mesos_scheduler.Event_MESSAGE:
		log.Print("Event MESSAGE: ", event)

	case mesos_scheduler.Event_FAILURE:
		log.Print("Event FAILURE: ", event)

	case mesos_scheduler.Event_ERROR:
		log.Print("Event ERROR: ", event)

	default:
		log.Print("Unexpected Event: ", event)
		return fmt.Errorf("Unexpected Event type: %q", event.Type)
	}

	return nil
}
