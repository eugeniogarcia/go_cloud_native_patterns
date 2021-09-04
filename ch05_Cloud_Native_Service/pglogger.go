package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"sync"

	_ "github.com/lib/pq" // Load the Postgres drivers
)

type PostgresDbParams struct {
	dbName   string
	host     string
	user     string
	password string
}

type PostgresTransactionLogger struct {
	events chan<- Event // Write-only channel for sending events
	errors <-chan error // Read-only channel for receiving errors
	db     *sql.DB      // Our database access interface
	wg     *sync.WaitGroup
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventPut, Key: key, Value: url.QueryEscape(value)}
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) LastSequence() uint64 {
	//No aplica. La base de datos gestiona esto
	return 0
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	go func() { // Query que ejecutaremos para insertar el evento en la base de datos
		query := `INSERT INTO transactions
			(event_type, key, value)
			VALUES ($1, $2, $3)`

		for e := range events { // Retrieve the next Event
			_, err := l.db.Exec( // Inserta el evento en la base de datos
				query,
				e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *PostgresTransactionLogger) Wait() {
	l.wg.Wait()
}

func (l *PostgresTransactionLogger) Close() error {
	l.wg.Wait()

	if l.events != nil {
		close(l.events) // Terminates Run loop and goroutine
	}

	return l.db.Close()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)    // An unbuffered events channel
	outError := make(chan error, 1) // A buffered errors channel

	//Query que ejecutaremos para leer los eventos de la base de datos
	query := "SELECT sequence, event_type, key, value FROM transactions"

	go func() {
		defer close(outEvent) // Close the channels when the
		defer close(outError) // goroutine ends

		//Ejecuta la query
		rows, err := l.db.Query(query) // Run query; get result set
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close() // This is important!

		var e Event // Create an empty Event

		//Recorre el cursor
		for rows.Next() { // Iterate over the rows

			//Lee el evento
			err = rows.Scan( // Read the values from the
				&e.Sequence, &e.EventType, // row into the Event.
				&e.Key, &e.Value)

			if err != nil {
				outError <- err
				return
			}

			//Publica el evento
			outEvent <- e // Send e to the channel
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

//Comprueba que la tabla exista
func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	const table = "transactions"

	var result string

	rows, err := l.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", table))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != table {
		rows.Scan(&result)
	}

	return result == table, rows.Err()
}

func (l *PostgresTransactionLogger) createTable() error {
	var err error

	//Tabla a crear
	createQuery := `CREATE TABLE transactions (
		sequence      BIGSERIAL PRIMARY KEY,
		event_type    SMALLINT,
		key 		  TEXT,
		value         TEXT
	  );`

	_, err = l.db.Exec(createQuery)
	if err != nil {
		return err
	}

	return nil
}

func NewPostgresTransactionLogger(param PostgresDbParams) (TransactionLogger, error) {
	//Cadena de conexión
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		param.host, param.dbName, param.user, param.password)

	//Nos conectamos
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create db value: %w", err)
	}

	//Chequea la conexión
	err = db.Ping() // Test the databases connection
	if err != nil {
		return nil, fmt.Errorf("failed to opendb connection: %w", err)
	}

	//Crea el logger
	tl := &PostgresTransactionLogger{db: db, wg: &sync.WaitGroup{}}

	//Comprueba que exista la tabla y sino la crea
	exists, err := tl.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err = tl.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return tl, nil
}
