package TinyActors

import (
	"github.com/pkg/errors"
	"time"
)

type State uint8
const(
	EndOfStream State = iota
)

const MessageMetadataReduce = "_reduce"
const MessageMetadataAnswer = "_answer"

type ActorModel struct {
	mailbox chan *Message
	system  *System
	run     func(actor *Actor)
}

func (system *System) Declare(action func(*Actor, *Message)) *ActorModel {
	typ := &ActorModel{
		make(chan *Message, 500),
		system,
		func(actor *Actor) {
			for {
				v, ok := <-actor.mailbox
				if ok {
					v.Receiver = actor
					action(actor, v)
				}
				if !ok || actor.dropped {
					break
				}
			}
		},
	}
	system.declare(typ)
	return typ
}

func (system *System) DeclareReducer(timeout time.Duration, size int, reduce func(*Actor, []*Message),
	action func(*Actor, *Message)) *ActorModel {

	typ := &ActorModel{
		make(chan *Message, 500),
		system,
		func(actor *Actor) {
			for {
				buff := make([]*Message, size)
				if len(buff) != size {
					panic("err")
				}
				for i := 0; i < size; i++ {
					select {
					case v, ok := <-actor.mailbox:
						if ok {
							if len(buff) != size {
								actor.mailbox <- v
							} else if buff != nil {
								buff[i] = v
							}
						}
					case <-time.After(timeout):
						for _, b := range buff {
							if b != nil {
								if b.Context[MessageMetadataReduce] == nil {
									b.Context[MessageMetadataReduce] = 1
									actor.mailbox <- b
								} else if b.Context[MessageMetadataReduce].(int) < 2 {
									b.Context[MessageMetadataReduce] = b.Context[MessageMetadataReduce].(int) + 1
									actor.mailbox <- b
								} else {
									action(actor, b)
								}
							}
						}
						buff = nil
					}
				}
				if buff != nil {
					reduce(actor, buff)
				}
			}
		},
	}
	system.declare(typ)
	return typ
}

func (typ *ActorModel) Tell(message *Message) {
	message.Sender = message.Receiver
	typ.mailbox <- message
}

func (typ *ActorModel) Forward(message *Message) {
	typ.mailbox <- message
}

func (typ *ActorModel) SimpleAsk(message *Message) chan *Message {
	answerChan := make(chan *Message)
	message.Context[MessageMetadataAnswer] = answerChan
	typ.mailbox <- message
	return answerChan
}

func (typ *ActorModel) Ask(message *Message, timeout time.Duration) (*Message, error) {
	answerChan := typ.SimpleAsk(message)

	select {
	case res := <-answerChan:
		return res, nil
	case <-time.After(timeout):
		return nil, errors.New("timed out")
	}
}

func (typ *ActorModel) Push(value interface{}) {
	typ.Tell(NewMessage(value))
}

func (typ *ActorModel) PushAsk(value interface{}, timeout time.Duration) (*Message, error) {
	return typ.Ask(NewMessage(value), timeout)
}

func (typ *ActorModel) instantiate() *Actor {
	return &Actor{
		typ,
		false,
	}
}

type Actor struct {
	*ActorModel
	dropped bool
}
