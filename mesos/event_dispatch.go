package mesos

import (
	"fmt"
	"log"

	//"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

func (m *MesosMaster) eventDispatch(event *mesos_scheduler.Event) error {

	switch *event.Type {
	case mesos_scheduler.Event_OFFERS:
		for _, offer := range event.Offers.Offers {
			if *offer.FrameworkId.Value != *m.frameworkId.Value {
				log.Print("Unexpected framework in offer: want %q, got %q",
					*m.frameworkId.Value, *offer.FrameworkId.Value)
				continue
			}

			if len(m.Offers) < cap(m.Offers) {
				m.Offers <- offer
			} else {
				// TODO(weingart): how to ignore/return offer?
			}
		}

	default:
		log.Print("Unexpected event: ", event)
		return fmt.Errorf("Unexpected evet: %+u", event)
	}

	return nil
}
