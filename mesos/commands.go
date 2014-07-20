package mesos

import (
	"fmt"
	"log"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

func (m *MesosMaster) LaunchTask(offer mesos.Offer, task *MesosTask) error {

	m.command <- func(fm *MesosMaster) error {

		launchType := mesos_scheduler.Call_LAUNCH
		launchCall := &mesos_scheduler.Call{
			FrameworkInfo: &mesos.FrameworkInfo{
				User: &m.config.RegisteredUser,
				Name: &m.config.FrameworkName,
				Id:   &m.frameworkId,
			},
			Type: &launchType,
			Launch: &mesos_scheduler.Call_Launch{
				TaskInfos: []*mesos.TaskInfo{
					&mesos.TaskInfo{
						Name: &task.Command,
						TaskId: &mesos.TaskID{
							Value: &task.Id,
						},
						SlaveId:   offer.SlaveId,
						Resources: offer.Resources,
						Command: &mesos.CommandInfo{
							Value: &task.Command,
						},
					},
				},
				OfferIds: []*mesos.OfferID{
					offer.Id,
				},
			},
		}

		return m.send(launchCall)
	}

	return nil
}

// XXX ---------------------------------------------------------------------
// XXX The "old way" is below
// XXX ---------------------------------------------------------------------

// Launch the given command as a task and consume the offer.
func (m *MesosMaster) LaunchTaskOld(offer mesos.Offer, id, command string) error {
	log.Printf("launching %s: %q", id, command)

	launchType := mesos_scheduler.Call_LAUNCH
	launchCall := &mesos_scheduler.Call{
		FrameworkInfo: &mesos.FrameworkInfo{
			User: &m.config.RegisteredUser,
			Name: &m.config.FrameworkName,
			Id:   &m.frameworkId,
		},
		Type: &launchType,
		Launch: &mesos_scheduler.Call_Launch{
			TaskInfos: []*mesos.TaskInfo{
				&mesos.TaskInfo{
					Name: &command,
					TaskId: &mesos.TaskID{
						Value: &id,
					},
					SlaveId:   offer.SlaveId,
					Resources: offer.Resources,
					Command: &mesos.CommandInfo{
						Value: &command,
					},
				},
			},
			OfferIds: []*mesos.OfferID{
				offer.Id,
			},
		},
	}

	err := m.send(launchCall)
	if err != nil {
		return err
	}

	// TODO(dhamon): Don't wait for the task to finish.
	for {
		event := <-m.events

		if *event.Type != mesos_scheduler.Event_UPDATE {
			return fmt.Errorf("unexpected event type: want %q, got %+v", mesos_scheduler.Event_UPDATE, *event.Type)
		}

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
			return nil
		case mesos.TaskState_TASK_RUNNING:
			log.Printf("task %s is running: %s", id, event.Update.Status.GetMessage())
			break
		case mesos.TaskState_TASK_FINISHED:
			log.Printf("task %s is complete: %s", id, event.Update.Status.GetMessage())
			return nil
		case mesos.TaskState_TASK_FAILED:
		case mesos.TaskState_TASK_KILLED:
		case mesos.TaskState_TASK_LOST:
			return fmt.Errorf("task %s failed to complete: %s", id, event.Update.Status.GetMessage())
		}
	}
	return nil
}
