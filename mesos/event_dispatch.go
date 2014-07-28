package mesos

import (
	"fmt"
	"log"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

// TODO(dhamon): custom driver logger

func (d *Driver) eventDispatch(event *mesos_scheduler.Event) error {
	switch *event.Type {
	case mesos_scheduler.Event_REGISTERED:
		log.Println("Event REGISTERED: ", event)

	case mesos_scheduler.Event_REREGISTERED:
		log.Println("Event REREGISTERED: ", event)

	case mesos_scheduler.Event_OFFERS:
		for _, offer := range event.Offers.Offers {
			if *offer.FrameworkId.Value != *d.frameworkId.Value {
				log.Printf("unexpected framework in offer: want %q, got %q",
					*d.frameworkId.Value, *offer.FrameworkId.Value)
				continue
			}

			if len(d.Offers) < cap(d.Offers) {
				d.Offers <- offer
			} else {
				// TODO(weingart): how to ignore/return offer?
				log.Printf("ignoring offer that we have no capacity for: %+v",
					offer)
			}
		}

	case mesos_scheduler.Event_RESCIND:
		log.Printf("Event RESCIND: %+v", event)

	case mesos_scheduler.Event_UPDATE:
		log.Printf("Event UPDATE: %+v", event)

		switch *event.Update.Status.State {
		case mesos.TaskState_TASK_STAGING,
			mesos.TaskState_TASK_STARTING,
			mesos.TaskState_TASK_RUNNING,
			mesos.TaskState_TASK_FINISHED,
			mesos.TaskState_TASK_FAILED,
			mesos.TaskState_TASK_KILLED,
			mesos.TaskState_TASK_LOST:

			d.Updates <- &TaskStateUpdate{
				TaskId:  event.Update.Status.GetTaskId().GetValue(),
				SlaveId: event.Update.Status.GetSlaveId().GetValue(),
				State:   event.Update.Status.GetState(),
				uuid:    event.Update.GetUuid(),
				driver:  d,
			}
		default:
			log.Printf("Unknown Event_UPDATE: %+v", event)
		}

	case mesos_scheduler.Event_MESSAGE:
		log.Printf("Event MESSAGE: %+v", event)

	case mesos_scheduler.Event_FAILURE:
		log.Printf("Event FAILURE: %+v", event)

	case mesos_scheduler.Event_ERROR:
		log.Printf("Event ERROR: %+v", event)

	default:
		log.Printf("Unexpected Event: %+v", event)
		return fmt.Errorf("unexpected event type: %q", event.Type)
	}

	return nil
}
