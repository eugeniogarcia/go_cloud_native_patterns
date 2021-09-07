package frontend

import (
	"egsmartin.com/hexarch/core"
)

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error
}

//Implementacion de FrontEnd
type zeroFrontEnd struct{}

func (f zeroFrontEnd) Start(kv *core.KeyValueStore) error {
	return nil
}
