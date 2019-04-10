package main

import (
	"fmt"
	ta "github.com/skyforce77/TinyActors"
	"os"
	_ "os"
	"time"
)

func main() {
	system := ta.NewSystem()

	actor1 := system.DeclareReducer(time.Second, 2, func(self *ta.Actor, message []*ta.Message) {
		self.Forward(message[0].New(
			message[0].Value.(int) + message[1].Value.(int),
		))
		time.Sleep(1000 * time.Millisecond)
	}, func(self *ta.Actor, message *ta.Message) {
		fmt.Println("=", message.Value)
		os.Exit(0)
	})

	system.Start()

	for i := 0; i < 500; i++ {
		actor1.Tell(1)
		time.Sleep(10 * time.Millisecond)
	}

	for{
		time.Sleep(time.Second)
	}
}
