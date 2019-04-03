package main

import (
	"fmt"
	ta "github.com/skyforce77/TinyActors"
	"time"
)

func main() {
	system := ta.NewSystem()

	actor2 := system.DeclareReducer(func(message1 *ta.Message, message2 *ta.Message) {
		fmt.Println(message.Value)
	})
	actor1 := system.Declare(func(message *ta.Message) {
		actor2.Forward(message)
		time.Sleep(time.Second)
	})

	system.Start()

	for i := 0; i < 500; i++ {
		actor1.Tell(i)
	}

	system.Finish()
}
