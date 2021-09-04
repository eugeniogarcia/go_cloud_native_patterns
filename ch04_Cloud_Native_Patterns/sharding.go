package ch04

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	sync.RWMutex
	m map[string]interface{}
}

type ShardedMap []*Shard

// NewShardedMap creates and initializes a new ShardedMap with the specified
// number of shards.
func NewShardedMap(nshards int) ShardedMap {
	//Crea el slice de Shards, con el tamaño y capacidad indicados
	shards := make([]*Shard, nshards)
	//Crea cada uno de los shards
	for i := 0; i < nshards; i++ {
		//Crea el mapa de cada Shard
		shard := make(map[string]interface{})
		shards[i] = &Shard{m: shard}
	}

	return shards
}

// getShardIndex accepts a key and returns a value in 0..N-1, where N is
// the number of shards. As currently written the hash algorithm only works
// correctly for up to 255 shards.
func (m ShardedMap) getShardIndex(key string) int {
	//convierte la key en un array de bytes, y calcula su hash. El hash es un array de 20 bytes
	hash := sha1.Sum([]byte(key))

	// Grab an arbitrary byte and mod it by the number of shards
	//Toma uno de los bytes y hace el mod para saber a que Shar corresponderá el key
	return int(hash[17]) % len(m)
}

// getShard accepts a key and returns a pointer to its corresponding Shard.
func (m ShardedMap) getShard(key string) *Shard {
	index := m.getShardIndex(key)
	return m[index]
}

// Delete removes a value from the map. If key doesn't exist in the map,
// this method is a no-op.
func (m ShardedMap) Delete(key string) {
	//Como modificaremos el contenido del Shard, adquitimos un lock de escritura - uno normal
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	delete(shard.m, key)
}

// Get retrieves and returns a value from the map. If the value doesn't exist,
// nil is returned.
func (m ShardedMap) Get(key string) interface{} {
	shard := m.getShard(key)
	//Adquirimos un lock de lectura
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

func (m ShardedMap) Set(key string, value interface{}) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = value
}

// Keys returns a list of all keys in the sharded map.
func (m ShardedMap) Keys() []string {
	//Crea un slice de strings, donde se incluirán todas las keys. El tamaño es cero
	keys := make([]string, 0)
	//Actualizaremos el slice anterior desde distintas go-rutinas, así que necesitaremos un mutex
	mutex := sync.Mutex{}

	wg := sync.WaitGroup{}
	//Crearemos una go-rutina para obtener los keys de cada Shard
	wg.Add(len(m))

	//Lanza las go-rutinas
	for _, shard := range m {
		go func(s *Shard) {
			//Bloqueamos el shar para lectura
			s.RLock()
			defer func() {
				s.RUnlock()
				wg.Done()
			}()

			//Recupera todas las keys del mapa asociado al Shard
			for key := range s.m {
				//Bloquea el slice...
				mutex.Lock()
				//... y lo actualiza
				keys = append(keys, key)
				mutex.Unlock()
			}

		}(shard)
	}

	//Esperamos a que hayan terminado todas las go-rutinas
	wg.Wait() // Block until all reads are done

	return keys // Return combined keys slice
}
