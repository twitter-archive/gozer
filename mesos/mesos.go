package mesos

import (
	"net"
	"os"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

func New(mc *MesosMasterConfig) (m *MesosMaster, err error) {

	name, err := os.Hostname()
	if err != nil {
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return
	}

	m = &MesosMaster{
		config:    *mc,
		command:   make(chan func(*MesosMaster) error),
		events:    make(chan *mesos_scheduler.Event, 100),
		Offers:    make(chan *mesos.Offer, 100),
		localIp:   addrs[0],
		localPort: 8888, // TODO(weingart): use ephemeral port
	}

	return
}
