package mesos

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"mesos.pb"
	"scheduler.pb"
)

var (
	master     = flag.String("master", "localhost", "Hostname of the master")
	masterPort = flag.Int("masterPort", 5050, "Port of the master")

	ip	string

	frameworkName string
	registeredUser string
	frameworkId mesos.FrameworkID

	httpWaitGroup sync.WaitGroup

	events = make(chan *mesos_scheduler.Event)
)

func init() {
	// Get our local IP address.
	name, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %+v", err)
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		log.Fatalf("Failed to get address for hostname %q: %+v", name, err)
	}

	log.Printf("using ip %q", addrs[0])
	ip = addrs[0]
}

// Register with a running master.
func Register(user, name string) error {
	frameworkName = name
	registeredUser = user

	// Create the register message and send it.
	callType := mesos_scheduler.Call_REGISTER
	registerCall := &mesos_scheduler.Call{
		Type: &callType,
		FrameworkInfo: &mesos.FrameworkInfo{
			User: &user,
			Name: &frameworkName,
		},
	}

	// Ensure we are listening before trying to send.
	httpWaitGroup.Wait()

	err := send(registerCall)
	if err != nil {
		return err
	}

	// Wait for HTTP endpoint to receive registration message.
	event := <-events

	if *event.Type != mesos_scheduler.Event_REGISTERED {
		return fmt.Errorf("Unexpected event type: want %q, got %+v", mesos_scheduler.Event_REGISTERED, *event.Type)
	}

	frameworkId = *event.Registered.FrameworkId
	log.Printf("Registered %s:%s with id %q", registeredUser, frameworkName, *frameworkId.Value)

	return nil
}

// Wait for offers.
func WaitForOffers() ([]mesos.Offer, error) {
	event := <-events

	if *event.Type != mesos_scheduler.Event_OFFERS {
		return nil, fmt.Errorf("Unexpected event type: want %q, got %+v", mesos_scheduler.Event_OFFERS, *event.Type)
	}

	var offers []mesos.Offer
	for _, offer := range event.Offers.Offers {
		if *offer.FrameworkId.Value != *frameworkId.Value {
			return nil, fmt.Errorf("Unexpected framework in offer: want %q, got %q", *frameworkId.Value, *offer.FrameworkId.Value)
		}
		offers = append(offers, *offer)
	}

	return offers, nil
}

// TODO(dhamon): pass in request types.
func RequestOffers() ([]mesos.Offer, error) {
	// Create the request message and send it.
	callType := mesos_scheduler.Call_REQUEST
	cpus := "cpus"
	memory := "memory"
	scalar := mesos.Value_SCALAR

	requestCall := &mesos_scheduler.Call {
		FrameworkInfo: &mesos.FrameworkInfo {
			User: &registeredUser,
			Name: &frameworkName,
			Id: &frameworkId,
		},
		Type: &callType,
		Request: &mesos_scheduler.Call_Request {
			Requests : []*mesos.Request {
				&mesos.Request {
					Resources: []*mesos.Resource {
						&mesos.Resource {
							Name: &cpus,
							Type: &scalar,
						},
						&mesos.Resource {
							Name: &memory,
							Type: &scalar,
						},
					},
				},
			},
		},
	}

	err := send(requestCall)
	if err != nil {
		return nil, err
	}

	return WaitForOffers()
}
