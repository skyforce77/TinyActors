package TinyActors

type Message struct {
	Value   interface{}
	Context map[string]interface{}
}

func (message *Message) New(value interface{}) *Message {
	return &Message{
		value,
		message.Context,
	}
}

func NewMessage(value interface{}) *Message {
	return &Message{
		value,
		make(map[string]interface{}),
	}
}
