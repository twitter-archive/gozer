package mesos

import (
	"fmt"
	"strings"

	"github.com/twitter/gozer/proto/mesos.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

type Offer struct {
	Id	   string
	driver	   *Driver
	mesosOffer *mesos.Offer
}

func (o *Offer) String() string {
	var resourceStr []string
	for _, resource := range o.mesosOffer.Resources {
		str := fmt.Sprintf("%s [%s]: ", *resource.Name, *resource.Role)
		switch *resource.Type {
		case mesos.Value_SCALAR:
			str += fmt.Sprintf("%.3f", *resource.Scalar.Value)

		case mesos.Value_SET:
			str += strings.Join(resource.Set.Item, ", ")

		case mesos.Value_RANGES:
			var rangeStr []string
			for _, value := range resource.Ranges.Range {
				rangeStr = append(rangeStr, fmt.Sprintf("[%d->%d]", *value.Begin, *value.End))
			}
			str += "[" + strings.Join(rangeStr, ", ") + "]"
		}

		resourceStr = append(resourceStr, str)
	}

	return fmt.Sprintf("%s: resources {%s} on slave %s",
		o.Id,
		strings.Join(resourceStr, ", "),
		*o.mesosOffer.SlaveId.Value)
}

func (o *Offer) Decline() {
	o.driver.command <- func(d *Driver) error {
		declineType := mesos_scheduler.Call_DECLINE
		declineCall := &mesos_scheduler.Call{
			FrameworkInfo: &mesos.FrameworkInfo{
				User: &d.config.RegisteredUser,
				Name: &d.config.FrameworkName,
				Id:   &d.frameworkId,
			},
			Type: &declineType,
			Decline: &mesos_scheduler.Call_Decline{
				OfferIds: []*mesos.OfferID{
					o.mesosOffer.Id,
				},
			},
		}

		return d.send(declineCall)
	}
}

