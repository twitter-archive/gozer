package mesos

import (
	"fmt"

	"code.google.com/p/goprotobuf/proto"

	"github.com/twitter/gozer/proto/messages.pb"
	"github.com/twitter/gozer/proto/scheduler.pb"
)

func bytesToEvent(protoType string, data []byte) (*mesos_scheduler.Event, error) {
	switch protoType {
	case "mesos.internal.FrameworkRegisteredMessage":
		message := new(mesos_internal.FrameworkRegisteredMessage)
		err := proto.Unmarshal(data, message)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %q into message of type %q: %+v", string(data), protoType, err)
		}
		eventType := mesos_scheduler.Event_REGISTERED
		return &mesos_scheduler.Event{
			Type: &eventType,
			Registered: &mesos_scheduler.Event_Registered{
				FrameworkId: message.FrameworkId,
				MasterInfo:  message.MasterInfo,
			},
		}, nil

	case "mesos.internal.FrameworkReregisteredMessage":
		message := new(mesos_internal.FrameworkReregisteredMessage)
		err := proto.Unmarshal(data, message)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %q into message of type %q: %+v", string(data), protoType, err)
		}
		eventType := mesos_scheduler.Event_REREGISTERED
		return &mesos_scheduler.Event{
			Type: &eventType,
			Reregistered: &mesos_scheduler.Event_Reregistered{
				FrameworkId: message.FrameworkId,
				MasterInfo:  message.MasterInfo,
			},
		}, nil

	case "mesos.internal.ResourceOffersMessage":
		message := new(mesos_internal.ResourceOffersMessage)
		err := proto.Unmarshal(data, message)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %q into message of type %q: %+v", string(data), protoType, err)
		}
		eventType := mesos_scheduler.Event_OFFERS
		return &mesos_scheduler.Event{
			Type: &eventType,
			Offers: &mesos_scheduler.Event_Offers{
				Offers: message.Offers,
			},
		}, nil

	case "mesos.internal.RescindResourceOfferMessage":
		message := new(mesos_internal.RescindResourceOfferMessage)
		err := proto.Unmarshal(data, message)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %q into message of type %q: %+v", string(data), protoType, err)
		}
		eventType := mesos_scheduler.Event_RESCIND
		return &mesos_scheduler.Event{
			Type: &eventType,
			Rescind: &mesos_scheduler.Event_Rescind{
				OfferId: message.OfferId,
			},
		}, nil

	case "mesos.internal.StatusUpdateMessage":
		message := new(mesos_internal.StatusUpdateMessage)
		err := proto.Unmarshal(data, message)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %q into message of type %q: %+v", string(data), protoType, err)
		}
		eventType := mesos_scheduler.Event_UPDATE
		return &mesos_scheduler.Event{
			Type: &eventType,
			Update: &mesos_scheduler.Event_Update{
				Uuid:   message.Update.Uuid,
				Status: message.Update.Status,
			},
		}, nil

	case "mesos.internal.ExitedExecutorMessage":
		message := new(mesos_internal.ExitedExecutorMessage)
		err := proto.Unmarshal(data, message)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %q into message of type %q: %+v", string(data), protoType, err)
		}
		eventType := mesos_scheduler.Event_FAILURE
		return &mesos_scheduler.Event{
			Type: &eventType,
			Failure: &mesos_scheduler.Event_Failure{
				SlaveId: message.SlaveId,
				ExecutorId: message.ExecutorId,
				Status: message.Status,
			},
		}, nil

	case "mesos.internal.ExecutorToFrameworkMessage":
		message := new(mesos_internal.ExecutorToFrameworkMessage)
		err := proto.Unmarshal(data, message)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %q into message of type %q: %+v", string(data), protoType, err)
		}
		eventType := mesos_scheduler.Event_MESSAGE
		return &mesos_scheduler.Event{
			Type: &eventType,
			Message: &mesos_scheduler.Event_Message{
				SlaveId: message.SlaveId,
				ExecutorId: message.ExecutorId,
				Data: data,
			},
		}, nil

	case "mesos.internal.FrameworkErrorMessage":
		message := new(mesos_internal.FrameworkErrorMessage)
		err := proto.Unmarshal(data, message)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %q into message of type %q: %+v", string(data), protoType, err)
		}
		eventType := mesos_scheduler.Event_ERROR
		return &mesos_scheduler.Event{
			Type: &eventType,
			Error: &mesos_scheduler.Event_Error{
				Message: message.Message,
			},
		}, nil
	}

	return nil, fmt.Errorf("unimplemented event type %q", protoType)
}
