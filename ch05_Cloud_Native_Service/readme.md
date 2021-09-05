# Ejecutar en container

## Crear la imagen

```ps
docker build --tag kvs:multipart .
```

Para ejecutarla:

```ps
docker run --rm --add-host="localhost:192.168.1.130" -d -p 84:8080 kvs:multipart
```

Para poder llamar desde el contenedor usando `localhost` y que se resuelva a la ip de nuestro host, `192.168.1.130`:

## Comentarios sobre el dockerfile

- Usamos un dockerfile __multistage_:
    - Ejecutamos los tests (`FROM golang:1.17 as test`)
    - Compilamos (`FROM golang:1.17 as build`)
    - Creamos el usuario - para no ejecutar la imagen como root (`FROM ubuntu:latest as user`)
    - Creamos la imagen  (`FROM scratch`)
- Al usar la imagen `golang` no hemos necesitado descargar nada relativo a go con _apt_. En _go.mod y go.sum_ tenemos la relacion de modulos de go que vamos a necesitar
- Al compilar tomamos los paquetes que se usaron para pasar - satisfactoriamente - los tests:
    
```txt
COPY --from=test /go/pkg/mod/ /go/pkg/mod/
```

- Creamos un usuario y le definimos como _owner_ del archivo que luego usaremos para guardar los transaction logs:

```txt
RUN useradd -u 10001 kv-user
RUN touch /transactions.log && chown kv-user /transactions.log
```

- Para que la imagen final sea lo más pequeña posible, usamos `` y copiamos de cada _stage_ lo que necesitamos - ejecutable, certificados, configuración de usuario, archivo del log de transacciones:

```txt
FROM scratch

# Copy the binary from the build container.
COPY --from=build /src/kvs .
COPY --from=build /src/*.pem .

COPY --from=user /etc/passwd /etc/passwd
COPY --from=user /transactions.log /transactions.log
```

## Configurar Postgress

Para poder aceptar conexiones desde 192.168.1.130 en mi postgress hemos creado la siguiente entrada en `pg_hba.conf`. Este archivo lo podemos localizar en `C:\Program Files\PostgreSQL\13\data`. La configuración inicial era:

```txt
# TYPE  DATABASE        USER            ADDRESS                 METHOD

# "local" is for Unix domain socket connections only
local   all             all                                     scram-sha-256
# IPv4 local connections:
host    all             all             127.0.0.1/32            scram-sha-256
# IPv6 local connections:
host    all             all             ::1/128                 scram-sha-256
# Allow replication connections from localhost, by a user with the
# replication privilege.
local   replication     all                                     scram-sha-256
host    replication     all             127.0.0.1/32            scram-sha-256
host    replication     all             ::1/128                 scram-sha-256
```

Hemos añadido una entrada:

```txt
# TYPE  DATABASE        USER            ADDRESS                 METHOD

# "local" is for Unix domain socket connections only
local   all             all                                     scram-sha-256
# IPv4 local connections:
host    all             all             127.0.0.1/32            scram-sha-256
host    all             all             192.168.1.130/32        scram-sha-256
# IPv6 local connections:
host    all             all             ::1/128                 scram-sha-256
# Allow replication connections from localhost, by a user with the
# replication privilege.
local   replication     all                                     scram-sha-256
host    replication     all             127.0.0.1/32            scram-sha-256
host    replication     all             ::1/128                 scram-sha-256
```
