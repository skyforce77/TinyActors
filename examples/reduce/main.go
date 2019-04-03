package main

import (
	"fmt"
	ta "github.com/skyforce77/TinyActors"
	"time"
)

func main() {
	system := ta.NewSystem()

	actor2 := system.DeclareReducer(2, func(self *ta.Actor, message []*ta.Message) {
		self.Forward(message[0].New(
			message[0].Value.(int) + message[1].Value.(int)),
		)
	}, func(message *ta.Message) {
		fmt.Println(message)
	})
	actor1 := system.Declare(func(message *ta.Message) {
		actor2.Forward(message)
		time.Sleep(time.Second)
	})

	system.Start()

	for i := 0; i < 500; i++ {
		actor1.Tell(1)
	}

	system.Finish()
}
