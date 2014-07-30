package mesos

import (
	"io/ioutil"
	"net"
	"os"

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
	config    DriverConfig
	pidIp     string
	pidPort   int

	frameworkId mesos.FrameworkID

	command chan func(*Driver) error
	// TODO(weingart): move to internal type to handle master disconnect, error events/etc.
	events chan *mesos_scheduler.Event

	Offers  chan *mesos.Offer
	Updates chan *TaskStateUpdate
}

func newDriver(mc *DriverConfig) (d *Driver, err error) {
	name, err := os.Hostname()
	if err != nil {
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return
	}

	d = &Driver{
		config:  *mc,
		pidIp:   addrs[0],
		pidPort: 8888, // TODO(weingart): use ephemeral port
		// TODO(dhamon): set channel filters
		log:	 NewLog("driver", ioutil.Discard, os.Stdout, os.Stdout, os.Stderr),
		command: make(chan func(*Driver) error),
		events:  make(chan *mesos_scheduler.Event, 100),
		Offers:  make(chan *mesos.Offer, 100),
		Updates: make(chan *TaskStateUpdate),
	}

	return
}

func New(framework, user, master string, port int) (d *Driver, err error) {
	cf := &DriverConfig{
		FrameworkName:  framework,
		RegisteredUser: user,
		Masters: []MasterAddress{
			MasterAddress{Hostname: master, Port: port},
		},
	}

	if d, err = newDriver(cf); err == nil {
		go d.Run()
	}

	return
}
