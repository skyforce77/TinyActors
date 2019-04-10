package main

import (
	"fmt"
	ta "github.com/skyforce77/TinyActors"
	"net/http"
	"time"
)

func main() {
	system := ta.NewSystem()

	actor2 := system.Declare(func(self *ta.Actor, message *ta.Message) {
		fmt.Println(message.Value)
	})
	actor1 := system.Declare(func(self *ta.Actor, message *ta.Message) {
		actor2.Forward(message)
		time.Sleep(time.Second)
	})

	system.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := ta.NewMessage(w)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
