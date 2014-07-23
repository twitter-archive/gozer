package mesos

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

type TaskStateUpdate struct {
	TaskId  string
	SlaveId string
	State   mesos.TaskState
	Uuid    uuid.UUID
	master  *MesosMaster
}

func (u *TaskStateUpdate) String() string {
	return fmt.Sprintf("update %s, task '%s' on slave '%s', state = %s",
		u.Uuid.String(),
		u.TaskId,
		u.SlaveId,
		u.State.String())
}

func (u *TaskStateUpdate) Ack() {

	u.master.command <- func(fm *MesosMaster) error {

		acknowledgeType := mesos_scheduler.Call_ACKNOWLEDGE
		acknowledgeCall := &mesos_scheduler.Call{
			FrameworkInfo: &mesos.FrameworkInfo{
				User: &fm.config.RegisteredUser,
				Name: &fm.config.FrameworkName,
				Id:   &fm.frameworkId,
			},
			Type: &acknowledgeType,
			Acknowledge: &mesos_scheduler.Call_Acknowledge{
				SlaveId: &mesos.SlaveID{
					Value: &u.SlaveId,
				},
				TaskId: &mesos.TaskID{
					Value: &u.TaskId,
				},
				Uuid: u.Uuid,
			},
		}

		return fm.send(acknowledgeCall)
	}
}
