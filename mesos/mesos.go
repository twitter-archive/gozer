package mesos

import (
	"fmt"
	"log"
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
		localIp:   addrs[0],
		localPort: 8888, // TODO(weingart): use ephemeral port
	}

	// Start up Framework http endpoint
	// TODO(weingart): move into state machine as part of m.Run()
	go startServing(m)

	return
}

// Register with a running master.
// TODO(weingart): make this an internal-only fuction
func (m *MesosMaster) Register(user, name string) error {

	// Create the register message and send it.
	callType := mesos_scheduler.Call_REGISTER
	registerCall := &mesos_scheduler.Call{
		Type: &callType,
		FrameworkInfo: &mesos.FrameworkInfo{
			User: &user,
			Name: &m.config.FrameworkName,
		},
	}

	err := m.send(registerCall)
	if err != nil {
		return err
	}

	// Wait for HTTP endpoint to receive registration message.
	event := <-m.events

	if *event.Type != mesos_scheduler.Event_REGISTERED {
		return fmt.Errorf("Unexpected event type: want %q, got %+v", mesos_scheduler.Event_REGISTERED, *event.Type)
	}

	m.frameworkId = *event.Registered.FrameworkId
	log.Printf("Registered %s:%s with id %q", m.config.RegisteredUser, m.config.FrameworkName, *m.frameworkId.Value)

	return nil
}

// TODO(dhamon): Refactor below into an event loop that tracks task updates and offers.
// Wait for offers.
func (m *MesosMaster) WaitForOffers() ([]mesos.Offer, error) {
	event := <-m.events

	if *event.Type != mesos_scheduler.Event_OFFERS {
		return nil, fmt.Errorf("Unexpected event type: want %q, got %+v", mesos_scheduler.Event_OFFERS, *event.Type)
	}

	var offers []mesos.Offer
	for _, offer := range event.Offers.Offers {
		if *offer.FrameworkId.Value != *m.frameworkId.Value {
			return nil, fmt.Errorf("Unexpected framework in offer: want %q, got %q", *m.frameworkId.Value, *offer.FrameworkId.Value)
		}
		offers = append(offers, *offer)
	}

	return offers, nil
}

// TODO(dhamon): pass in request types.
func (m *MesosMaster) RequestOffers() ([]mesos.Offer, error) {
	// Create the request message and send it.
	callType := mesos_scheduler.Call_REQUEST
	cpus := "cpus"
	memory := "memory"
	scalar := mesos.Value_SCALAR

	requestCall := &mesos_scheduler.Call{
		FrameworkInfo: &mesos.FrameworkInfo{
			User: &m.config.RegisteredUser,
			Name: &m.config.FrameworkName,
			Id:   &m.frameworkId,
		},
		Type: &callType,
		Request: &mesos_scheduler.Call_Request{
			Requests: []*mesos.Request{
				&mesos.Request{
					Resources: []*mesos.Resource{
						&mesos.Resource{
							Name: &cpus,
							Type: &scalar,
						},
						&mesos.Resource{
							Name: &memory,
							Type: &scalar,
						},
					},
				},
			},
		},
	}

	err := m.send(requestCall)
	if err != nil {
		return nil, err
	}

	return m.WaitForOffers()
}
