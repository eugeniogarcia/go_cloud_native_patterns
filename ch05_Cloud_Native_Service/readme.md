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
