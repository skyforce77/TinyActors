package TinyActors

type ActorType struct {
	mailbox chan *Message
	action func(*Message)
	system *System
}

func (system *System) Declare(action func(*Message)) *ActorType {
	typ := &ActorType{
		make(chan *Message, 500),
		action,
		system,
	}
	system.declare(typ)
	return typ
}

func (typ *ActorType) Forward(message *Message) {
	typ.mailbox <- message
}

func (typ *ActorType) Tell(value interface{}) {
	message := newMessage(value)
	typ.mailbox <- message
}

func (typ *ActorType) instanciate() *Actor {
	return &Actor{
		typ,
		false,
	}
}

type Actor struct {
	*ActorType
	dropped bool
}

func (actor *Actor) run() {
	for {
		v, ok := <-actor.mailbox
		if ok {
			actor.action(v)
		}
		if !ok || actor.dropped {
			break
		}
	}
}
