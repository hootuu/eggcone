package eggmq

type Payload interface {
	Of(str string)
	To() string
}
