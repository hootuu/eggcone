package eggmq

const (
	INIT       = 0
	PROCESSING = 1
	PENDING    = 2
	COMPLETED  = 3
	FAILED     = 4
)

type Message struct {
	ID         string `json:"id"`
	Topic      string `json:"topic"`
	Payload    string `json:"payload"`
	RetryCount int    `json:"retry_count"`
	Status     int    `json:"status"`
}

func NewMessage(id string, topic string, payload string) *Message {
	return &Message{
		ID:         id,
		Topic:      topic,
		Payload:    payload,
		RetryCount: 0,
		Status:     INIT,
	}
}
