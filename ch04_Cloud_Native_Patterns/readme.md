# Stability Patterns

## circuit breaker

```go
func Breaker(circuit Circuit, failureThreshold uint) Circuit {
```

Crea un _Breaker_ que llama a un _Circuit_ admitiendo _failureThreshold_ fallos. Si se producen _failureThreshold_ fallos consecutivos, el Breaker se abre de modo que se filtran las llamadas al circuito. Pasados _2^fallos segs_ se vuelve a intentar.

## DebounceFirst

```go
func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
```

Llamamos a un circuito. Si han pasado d segundos desde la ultima vez que usamos `DebounceFirst`, se llama al circuito y cachea la respuesta, sino devolvemos la respuesta cacheada.

## DebounceLast

```go
func DebounceLast(circuit Circuit, d time.Duration) Circuit {
```

Llamamos a un circuito. Hasta que no hayan pasado d segundos desde la ultima vez que usamos `DebounceLast`, no se hace la llamada al circuito, y se devuelve en su lugar la respuesta que tenemos cacheada - la que corresponde a la última vez que hicimos una llamada.

## Retry

```go
func Retry(effector Effector, retries int, delay time.Duration) Effector {
```

Llama a un _effector_ y si se produce un error, lo reuntenta _retries_ veces, espaciadas _delay_ segundos.

## Throttle

```go
func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
```

Limitamos el número de ejecuciones:
- _max_. Número máximo de ejecuciones
- _d_. Frecuencia con la que se incrementa la cuota
- _refill_. Valor en el que se incrementa la cuota. La cuota no podrá superar nunca el valor _max_

# Concurrency Patterns

## Future

Con el tipo _future_ definimos un tipo que usaremos para aquellos casos en los que necesitamos una ejecución asíncrona. 

```go
type Future interface {
	Result() (string, error)
}

type InnerFuture struct {
	once sync.Once
	wg   sync.WaitGroup

	res   string
	err   error
	resCh <-chan string
	errCh <-chan error
}
```

Llamando a _Result_ se bloqueara la ejecución hasta que la respuesta este lista:

```go
func (f *InnerFuture) Result() (string, error) {
	//Se ejecuta una vez, y se bloquea la ejecución hasta no recibir algo por los canales
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()
		f.res = <-f.resCh
		f.err = <-f.errCh
	})

	//Bloquea hasta que once termine
	f.wg.Wait()

	return f.res, f.err
}
```

## Sharding

Se trata de definir una estructura de datos que soporte el sharding en varios mapas. Un Shard tiene el mapa y un RW Mutex:

```go
type Shard struct {
	sync.RWMutex
	m map[string]interface{}
}
```

El otro tipo que necesitaremos es el mapa de Shards. Guarda un slice de Shards:

```go
type ShardedMap []*Shard
```

Accederemos a los datos a traves de `ShardedMap`:
- NewShardedMap. Crea el ShardMap
- `func (m ShardedMap) getShardIndex(key string) int {` distribuye cada Key en su shard, usando una técnica de hashing
- Acceso a datos:
    - `func (m ShardedMap) getShard(key string) *Shard {`. Obtiene el Shard que corresponde a una key
    - `func (m ShardedMap) Delete(key string) {`. Borra una entrada 
    - `func (m ShardedMap) Get(key string) interface{} {`.  Obtiene una entrada
    - `func (m ShardedMap) Set(key string, value interface{}) {`. Crea una entrada
    - `func (m ShardedMap) Keys() []string {`. Lista todas las entradas
