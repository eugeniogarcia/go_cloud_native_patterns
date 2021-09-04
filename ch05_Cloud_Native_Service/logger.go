package main

type EventType byte

//Declara dos constantes de tipo EventType
const (
	_                     = iota // iota == 0; ignore this value
	EventDelete EventType = iota // iota == 1
	EventPut                     // iota == 2; implicitly repeat last
)

//Estructura de un evento
type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	LastSequence() uint64

	Run()
	Wait()
	Close() error

	ReadEvents() (<-chan Event, <-chan error)
}
