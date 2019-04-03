package main

import (
	"fmt"
	ta "github.com/skyforce77/TinyActors"
	"time"
)

func main() {
	system := ta.NewSystem()

	actor1 := system.DeclareReducer(2, func(self *ta.Actor, message []*ta.Message) {
		fmt.Println(message[0].Value.(int), message[1].Value.(int), message[0].Value.(int) + message[1].Value.(int))
		self.Forward(message[0].New(
			message[0].Value.(int) + message[1].Value.(int)),
		)
	}, func(self *ta.Actor, message *ta.Message) {
		fmt.Println(message)
	})

	system.Start()

	for i := 0; i < 500; i++ {
		actor1.Tell(1)
	}

	for{
		time.Sleep(100000)
	}
}
