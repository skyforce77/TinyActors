package main

import (
	"fmt"
	ta "github.com/skyforce77/TinyActors"
	"time"
)

type AddJob struct {
	a int
	b int
}

func main() {
	system := ta.NewSystem()

	actor1 := system.Declare(func(self *ta.Actor, message *ta.Message) {
		if add, ok := message.Value.(AddJob); ok {
			message.Answer(ta.NewMessage(add.a + add.b))
		}
	})

	system.Start()
	defer system.Finish()

	result, err := actor1.PushAsk(AddJob{10, 5}, 10*time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result: %d", result.Value.(int))
}
