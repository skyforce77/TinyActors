package TinyActors

type Message struct {
	Value    interface{}
	Context  map[string]interface{}
	Sender   *Actor
	Receiver *Actor
}

func (message *Message) New(value interface{}) *Message {
	return &Message{
		value,
		message.Context,
		message.Sender,
		message.Receiver,
	}
}

func NewMessage(value interface{}) *Message {
	return &Message{
		value,
		make(map[string]interface{}),
		nil,
		nil,
	}
}

func (message *Message) IsWaitingForAnAnswer() bool {
	if _, ok := message.Context[MessageMetadataAnswer]; ok {
		return true
	}
	return false
}

func (message *Message) Answer(answer *Message) bool {
	if val, ok := message.Context[MessageMetadataAnswer]; ok {
		if c, ok := val.(chan *Message); ok {
			c <- answer
		}
	}
	return false
}
