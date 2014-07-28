package mesos

import (
	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

type MasterAddress struct {
	Hostname string
	Port     int
}

type DriverConfig struct {
	FrameworkName  string
	RegisteredUser string
	Masters        []MasterAddress
}

type Driver struct {
	config  DriverConfig
	pidIp   string
	pidPort int

	frameworkId mesos.FrameworkID

	command chan func(*Driver) error
	// TODO(weingart): move to internal type to handle master disconnect, error events/etc.
	events chan *mesos_scheduler.Event

	Offers  chan *mesos.Offer
	Updates chan *TaskStateUpdate
}

type MesosTask struct {
	Id      string
	Command string
}
