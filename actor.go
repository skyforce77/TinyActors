package TinyActors

type ActorModel struct {
	mailbox chan *Message
	system  *System
	run func(actor *Actor)
}

func (system *System) Declare(action func(*Message)) *ActorModel {
	typ := &ActorModel{
		make(chan *Message, 500),
		system,
		func(actor *Actor) {
			for {
				v, ok := <-actor.mailbox
				if ok {
					action(v)
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

func (system *System) DeclareReducer(size int, action func([]*Message)) *ActorModel {
	typ := &ActorModel{
		make(chan *Message, 500),
		system,
		func(actor *Actor) {
			for {
				buff := make([]*Message, size)
				for i:=0; i<size; i++ {
					v, ok := <-actor.mailbox
					if ok {
						buff[i] = v
					}
					if !ok || actor.dropped {
						return
					}
				}
				action(buff)
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
