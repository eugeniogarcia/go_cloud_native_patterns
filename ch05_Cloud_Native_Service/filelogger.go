package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sync"
)

//Logger
type FileTransactionLogger struct {
	events       chan<- Event // Write-only channel for sending events
	errors       <-chan error
	lastSequence uint64   // The last used event sequence number
	file         *os.File // The location of the transaction log
	wg           *sync.WaitGroup
}

//Publica un evento PUT. Esta pensado para ser llamado como go-rutina
func (l *FileTransactionLogger) WritePut(key, value string) {
	l.wg.Add(1)
	//Notese que no se informa la secuencia en el evento
	l.events <- Event{EventType: EventPut, Key: key, Value: url.QueryEscape(value)}
}

//Publica un evento DELETE. Esta pensado para ser llamado como go-rutina
func (l *FileTransactionLogger) WriteDelete(key string) {
	l.wg.Add(1)
	//Notese que no se informa la secuencia en el evento
	l.events <- Event{EventType: EventDelete, Key: key}
}

//Obtenemos el canal de errores. El canal es un elemto privado, por eso necesitamos este método
func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

//Recuperamos la última secuencia usada
func (l *FileTransactionLogger) LastSequence() uint64 {
	return l.lastSequence
}

//Crea un logger con el archivo indicado. No se crean los canales. Para ello hay que llamar a Run
func NewFileTransactionLogger(filename string) (*FileTransactionLogger, error) {
	// Open the transaction log file for reading and writing.
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}
	return &FileTransactionLogger{file: file, wg: &sync.WaitGroup{}}, nil
}

func (l *FileTransactionLogger) Run() {
	//Crea el canal de eventos, que podra encolar hasta 16 eventos
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	// Start retrieving events from the events channel and writing them
	// to the transaction log
	go func() {
		//Con cada evento recibido
		for e := range events {
			//incrementamos la secuencia
			l.lastSequence++

			//lo escribimos en el log, incluyendo la secuencia
			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- fmt.Errorf("cannot write to log file: %w", err)
			}

			l.wg.Done()
		}
	}()
}

//Esperamos a que terminen las go-rutinas que puedan estar escribiendo eventos
func (l *FileTransactionLogger) Wait() {
	l.wg.Wait()
}

//Cierra el canal de eventos y el archivo
func (l *FileTransactionLogger) Close() error {
	//Esperamos a que terminen las go-rutinas que puedan estar escribiendo eventos
	l.wg.Wait()

	if l.events != nil {
		close(l.events) // Terminates Run loop and goroutine
	}

	return l.file.Close()
}

//Lee de un archivo con el formato de log las entradas y los publica en un canal de eventos que el propio método crea - y que no tiene que nada que ver con los canales del logger - en forma de Event
//Esta diseñado para ser usado en la inicialización
func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	//Busca cambios en el archivo
	scanner := bufio.NewScanner(l.file)
	//Crea los canales de salida
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	//Procesa de forma asincrona los cambios en el archivo
	go func() {
		//Evento
		var e Event

		defer close(outEvent)
		defer close(outError)

		//lee el contenido del archivo
		for scanner.Scan() {
			line := scanner.Text()

			//lee la nueva entrada
			fmt.Sscanf(
				line, "%d\t%d\t%s\t%s",
				&e.Sequence, &e.EventType, &e.Key, &e.Value) //actualiza el evento con los datos leidos

			//Comprueba que la secuencia que hemos leido sea mayor a la última procesada - monotonico
			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			//Procesa el valor leido, que tenga el formato correcto
			uv, err := url.QueryUnescape(e.Value)
			if err != nil {
				outError <- fmt.Errorf("vaalue decoding failure: %w", err)
				return
			}

			//Actualiza el evento de nuevo, especificamente el valor leido
			e.Value = uv
			//Guardamos la secuencia que hemos procesado
			l.lastSequence = e.Sequence

			//publicamos el evento
			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}
