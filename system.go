package TinyActors

import (
	"time"
)

type SystemState uint8
const(
	Created SystemState = iota
	Started
	Finishing
	Finished
)

type System struct {
	actorTypes []*ActorType
	actors map[*ActorType][]*Actor
	state SystemState
}

func NewSystem() *System {
	return &System{
		make([]*ActorType, 0),
		make(map[*ActorType][]*Actor),
		Created,
	}
}

func (system *System) declare(typ *ActorType) {
	system.actorTypes = append(system.actorTypes, typ)
}

func (system *System) addActor(typ *ActorType) {
	if system.actors[typ] == nil {
		system.actors[typ] = make([]*Actor, 0)
	}
	actor := typ.instanciate()
	system.actors[typ] = append(system.actors[typ], actor)
	go actor.run()
}

func (system *System) dropActor(typ *ActorType) {
	if system.actors[typ] == nil || len(system.actors[typ]) < 1 {
		return
	}
	system.actors[typ][0].dropped = true
	system.actors[typ] = system.actors[typ][1:]
}

func (system *System) monitor() {
	for {
		for _, typ := range system.actorTypes {
			if len(typ.mailbox) > cap(typ.mailbox)/4 || len(typ.mailbox) > 0 && len(system.actors[typ]) == 0 {
				system.addActor(typ)
			} else if len(system.actors[typ]) > 1 && len(system.actors[typ]) > 1 {
				if len(typ.mailbox) < cap(typ.mailbox)/10 {
					system.dropActor(typ)
				}
			}

			if system.state == Finishing && len(typ.mailbox) == 0 {
				system.dropActor(typ)
			}

			i := 0
			for _, acts := range system.actors {
				i += len(acts)
			}
			if i == 0 {
				system.state = Finished
			}
		}

		time.Sleep(20 * time.Millisecond)
	}
}

func (system *System) Start() {
	system.state = Started
	for _, typ := range system.actorTypes {
		system.addActor(typ)
	}
	go system.monitor()
}

func (system *System) Stop() {
	for _, typ := range system.actorTypes {
		for _, act := range system.actors[typ] {
			act.dropped = true
		}
	}
}

func (system *System) Finish() {
	system.state = Finishing

	for system.state != Finished {
		time.Sleep(20 * time.Millisecond)
	}
}