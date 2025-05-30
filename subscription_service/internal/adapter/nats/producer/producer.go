package producer

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventpb "github.com/recktt77/projectProto-definitions/gen/auth_service/genproto/subscription"
)

type Producer struct {
	nc *nats.Conn
}

func NewProducer(nc *nats.Conn) *Producer {
	return &Producer{nc: nc}
}

func (p *Producer) PublishEvent(subject string, eventType string, payload any) error {
	protoPayload, ok := payload.(proto.Message)
	if !ok {
		return fmt.Errorf("payload does not implement proto.Message")
	}
	anyPayload, err := anypb.New(protoPayload)

	if err != nil {
		return err
	}

	e := &eventpb.Event{
		Type: eventType,
		Data: anyPayload,
		Time: timestamppb.New(time.Now()),
	}

	data, err := proto.Marshal(e)
	if err != nil {
		return err
	}

	return p.nc.Publish(subject, data)
}
