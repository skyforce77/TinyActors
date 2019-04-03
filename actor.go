package TinyActors

import (
	"fmt"
	"time"
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

func (system *System) DeclareReducer(size int, reduce func(*Actor, []*Message),
	action func(*Actor, *Message)) *ActorModel {

	typ := &ActorModel{
		make(chan *Message, 500),
		system,
		func(actor *Actor) {
			for {
				buff := make([]*Message, size)
				for i := 0; i < size; i++ {
					select {
					case v, ok := <-actor.mailbox:
						if ok {
							buff[i] = v
						}
						if !ok || actor.dropped {
							buff = nil
						}
					case <-time.After(100 * time.Millisecond):
						for _, b := range buff {
							fmt.Println(buff)
							if b != nil {
								action(actor, b)
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
	message := newMessage(value)
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
