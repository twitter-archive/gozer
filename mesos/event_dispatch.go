package mesos

import (
	"fmt"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

func (d *Driver) eventDispatch(event *mesos_scheduler.Event) error {
	switch *event.Type {
	case mesos_scheduler.Event_REGISTERED:
		d.log.Info.Println("Event REGISTERED:", event)

	case mesos_scheduler.Event_REREGISTERED:
		d.log.Info.Println("Event REREGISTERED:", event)

	case mesos_scheduler.Event_OFFERS:
		for _, offer := range event.Offers.Offers {
			if *offer.FrameworkId.Value != *d.frameworkId.Value {
				d.log.Warn.Printf("unexpected framework in offer: want %q, got %q",
					*d.frameworkId.Value, *offer.FrameworkId.Value)
				continue
			}

			if len(d.Offers) < cap(d.Offers) {
				d.Offers <- offer
			} else {
				// TODO(weingart): how to ignore/return offer?
				d.log.Warn.Println("ignoring offer that we have no capacity for:", offer)
			}
		}

	case mesos_scheduler.Event_RESCIND:
		d.log.Info.Printf("Event RESCIND: %+v", event)

	case mesos_scheduler.Event_UPDATE:
		d.log.Info.Printf("Event UPDATE: %+v", event)

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
			d.log.Error.Printf("Unknown Event_UPDATE: %+v", event)
		}

	case mesos_scheduler.Event_MESSAGE:
		d.log.Info.Printf("Event MESSAGE: %+v", event)

	case mesos_scheduler.Event_FAILURE:
		d.log.Info.Printf("Event FAILURE: %+v", event)

	case mesos_scheduler.Event_ERROR:
		d.log.Info.Printf("Event ERROR: %+v", event)

	default:
		err := fmt.Errorf("unexpected event type: %q", event.Type)
		d.log.Error.Println(err)
		return err
	}

	return nil
}
