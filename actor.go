package TinyActors

import (
	"time"
)

type State uint8
const(
	EndOfStream State = iota
)

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
								if b.Context["_reduce"] == nil {
									b.Context["_reduce"] = 1
									actor.mailbox <- b
								} else if b.Context["_reduce"].(int) < 2 {
									b.Context["_reduce"] = b.Context["_reduce"].(int) + 1
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

func (typ *ActorModel) Forward(message *Message) {
	typ.mailbox <- message
}

func (typ *ActorModel) Tell(value interface{}) {
	message := NewMessage(value)
	typ.mailbox <- message
}

func (typ *ActorModel) instanciate() *Actor {
	return &Actor{
		typ,
		false,
	}
}

type Actor struct {
	*ActorModel
	dropped bool
}
