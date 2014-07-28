package mesos

import (
	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

func (d *Driver) LaunchTask(offer *mesos.Offer, task *MesosTask) error {
	d.command <- func(fm *Driver) error {
		launchType := mesos_scheduler.Call_LAUNCH
		launchCall := &mesos_scheduler.Call{
			FrameworkInfo: &mesos.FrameworkInfo{
				User: &fm.config.RegisteredUser,
				Name: &fm.config.FrameworkName,
				Id:   &fm.frameworkId,
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

		return fm.send(launchCall)
	}

	return nil
}
