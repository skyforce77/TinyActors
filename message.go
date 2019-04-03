package TinyActors

type Message struct {
	Value interface{}
	Context map[string]interface{}
}

func (message *Message) New(value interface{}) *Message {
	return &Message{
		value,
		message.Context,
	}
}

func newMessage(value interface{}) *Message {
	return &Message{
		value,
		make(map[string]interface{}),
	}
}
