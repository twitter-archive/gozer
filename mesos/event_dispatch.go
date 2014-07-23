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

		switch *event.Update.Status.State {
		case mesos.TaskState_TASK_STAGING,
			mesos.TaskState_TASK_STARTING,
			mesos.TaskState_TASK_RUNNING,
			mesos.TaskState_TASK_FINISHED,
			mesos.TaskState_TASK_FAILED,
			mesos.TaskState_TASK_KILLED,
			mesos.TaskState_TASK_LOST:

			m.Updates <- &TaskStateUpdate{
				TaskId:  event.Update.Status.GetTaskId().GetValue(),
				SlaveId: event.Update.Status.GetSlaveId().GetValue(),
				State:   event.Update.Status.GetState(),
				Uuid:    event.Update.GetUuid(),
				master:  m,
			}
		default:
			log.Print("Unknown Event_UPDATE: ", event)
		}

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
