package mesos

import (
	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

type MesosMasterLocation struct {
	Hostname string
	Port     int
}

type MesosMasterConfig struct {
	FrameworkName  string
	RegisteredUser string
	Masters        []MesosMasterLocation
}

type Driver struct {
	config    MesosMasterConfig
	localIp   string
	localPort int

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
