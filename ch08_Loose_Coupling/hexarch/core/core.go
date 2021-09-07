package core

import (
	"errors"
	"log"
	"sync"
)

type KeyValueStore struct {
	sync.RWMutex
	m        map[string]string
	transact TransactionLogger
}

var ErrorNoSuchKey = errors.New("no such key")

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		m:        make(map[string]string),
		transact: ZeroTransactionLogger{},
	}
}

func (store *KeyValueStore) Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	//Publica el evento para que se actualice el transaction log
	store.transact.WriteDelete(key)

	return nil
}

func (store *KeyValueStore) Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func (store *KeyValueStore) Put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	//Publica el evento para que se actualice el transaction log
	store.transact.WritePut(key, value)

	return nil
}

func (store *KeyValueStore) WithTransactionLogger(tl TransactionLogger) *KeyValueStore {
	store.transact = tl
	return store
}

//Leemos desde el transaction log los datos y los cargamos en el mapa.
func (store *KeyValueStore) Restore() error {
	var err error

	//Lee los datos y los publica en estos dos canales
	events, errors := store.transact.ReadEvents()
	count, ok, e := 0, true, Event{}

	//Procesa los eventos que se publican
	for ok && err == nil {
		select {
		case err, ok = <-errors:

		//Cuando el canal se cierre, ok será true, y se termina el while loop
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Got a DELETE event!
				err = store.Delete(e.Key)
				count++
			case EventPut: // Got a PUT event!
				err = store.Put(e.Key, e.Value)
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)

	//Lanzamos una gorutina que escuchara los eventos que se publiquen
	store.transact.Run()

	go func() {
		for err := range store.transact.Err() {
			log.Print(err)
		}
	}()

	return err
}

type ZeroTransactionLogger struct{}

func (z ZeroTransactionLogger) WriteDelete(key string)                   {}
func (z ZeroTransactionLogger) WritePut(key, value string)               {}
func (z ZeroTransactionLogger) Err() <-chan error                        { return nil }
func (z ZeroTransactionLogger) LastSequence() uint64                     { return 0 }
func (z ZeroTransactionLogger) Run()                                     {}
func (z ZeroTransactionLogger) Wait()                                    {}
func (z ZeroTransactionLogger) Close() error                             { return nil }
func (z ZeroTransactionLogger) ReadEvents() (<-chan Event, <-chan error) { return nil, nil }
