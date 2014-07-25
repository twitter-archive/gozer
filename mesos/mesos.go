package mesos

import (
	"log"
	"net"
	"os"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

func NewDriver(mc *DriverConfig) (d *Driver, err error) {

	name, err := os.Hostname()
	if err != nil {
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return
	}
	log.Printf("XXX: %s = %+v", name, addrs)

	d = &Driver{
		config:    *mc,
		command:   make(chan func(*Driver) error),
		events:    make(chan *mesos_scheduler.Event, 100),
		Offers:    make(chan *mesos.Offer, 100),
		Updates:   make(chan *TaskStateUpdate),
		localIp:   addrs[0],
		localPort: 8888, // TODO(weingart): use ephemeral port
	}

	return
}

func New(framework string, user string, master string, port int) (d *Driver, err error) {

	cf := &DriverConfig{
		FrameworkName:  framework,
		RegisteredUser: user,
		Masters: []MasterAddress{
			MasterAddress{Hostname: master, Port: port},
		},
	}

	return NewDriver(cf)
}
