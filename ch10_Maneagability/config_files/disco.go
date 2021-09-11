package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-yaml/yaml"
)

//Usamos el hash para comprobar si el contenido de un archivo ha cambiado. Esta función calcula el hash
func calculaHash(filepath string) (string, error) {
	file, err := os.Open(filepath) // Open the file for reading
	if err != nil {
		return "", err
	}
	defer file.Close() // Be sure to close your file!

	hash := sha256.New() // Use the Hash in crypto/sha256
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", hash.Sum(nil)) // Get encoded hash sum
	return sum, nil
}

//Comprueba períodicamente si el hash del archivo ha cambiado
func vigilaConfig(ctx context.Context, filepath string) (<-chan string, <-chan error, error) {
	//Crea los canales
	errs := make(chan error)
	changes := make(chan string)
	//Guarda el último valor del hash
	hash := ""
	//Obtiene el hash períodicamente
	go func() {
		ticker := time.NewTicker(time.Second)
		continua := true

		for continua {
			select {
			case <-ctx.Done():
				close(errs)
				close(changes)
				continua = false
			case <-ticker.C:
				//Calculamos el hash
				newhash, err := calculaHash(filepath)
				if err != nil {
					errs <- err
					continue
				}
				if hash != newhash {
					hash = newhash
					changes <- filepath
				}
			}
		}
	}()

	return changes, errs, nil
}

//Equivalente a vigilaConfig, pero usando un watcher del SSOO
func watchConfigNotify(ctx context.Context, filepath string) (<-chan string, <-chan error, error) {
	//Crea el canal
	changes := make(chan string)
	//Crea un watcher
	watcher, err := fsnotify.NewWatcher() // Get an fsnotify.Watcher
	if err != nil {
		return nil, nil, err
	}
	//Podemos añadir uno o varios archivos a vigilar
	err = watcher.Add(filepath) // Tell watcher to watch
	if err != nil {             // our config file
		return nil, nil, err
	}

	//Monitoriza los cambios
	go func() {
		defer watcher.Close()
		changes <- filepath // First is ALWAYS a change
		continua := true
		for continua {
			select {
			case <-ctx.Done():
				close(changes)
				continua = false
			case event := <-watcher.Events: // Range over watcher events
				//Se ha producido un evento
				//Comprobamos si el evento nos informa de una escritura
				if event.Op&fsnotify.Write == fsnotify.Write {
					//El nombre nos indica que archivo cambio
					changes <- event.Name
				}
			}
		}
	}()

	return changes, watcher.Errors, nil
}

func procesaCambios(updates <-chan string, errors <-chan error) {
	var filepath string
	var ok bool = true
	var err error

	for ok {
		select {
		case filepath, ok = <-updates:
			c, err := cargaConfiguracion(filepath)
			if err != nil {
				log.Println("error loading config:", err)
				continue
			}
			config = c
		case err, ok = <-errors:
			log.Println("error watching config:", err)

		}
	}
}

func cargaConfiguracion(filepath string) (Config, error) {
	dat, err := ioutil.ReadFile(filepath) // Ingest file as []byte
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = yaml.Unmarshal(dat, &config) // Do the unmarshal
	if err != nil {
		return Config{}, err
	}
	config.imprime()
	return config, nil
}

//Funcion que cancela el monitoreo del archivo de configuración
var cancela context.CancelFunc

//Configuración
var config Config

//Tenemos dos formas de vigilar los cambios en la configuracion. Usando un watcher del SSOO o haciendolo custom
var usa_watcher = true

func init() {
	var ctx context.Context
	var updates <-chan string
	var errors <-chan error
	var err error

	ctx, cancela = context.WithCancel(context.Background())

	if usa_watcher {
		updates, errors, err = watchConfigNotify(ctx, "config.yaml")
	} else {
		updates, errors, err = vigilaConfig(ctx, "config.yaml")

	}
	if err != nil {
		panic(err)
	}

	go procesaCambios(updates, errors)
}
