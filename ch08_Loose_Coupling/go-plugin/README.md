## Build shared object (`.so`) files.

```bash
go build -buildmode=plugin -o duck/duck.so duck/duck.go
go build -buildmode=plugin -o frog/frog.so frog/frog.go
go build -buildmode=plugin -o fox/fox.so fox/fox.go
```

## Check file types

Just for fun, check the file types.

```bash
file duck/duck.so
file frog/frog.so
file fox/fox.so
```

Depending on your OS, you'll see something like:

```bash
file duck/duck.so

duck/duck.so: ELF 64-bit LSB shared object, x86-64, version 1 (SYSV), dynamically linked, BuildID[sha1]=82e82042bafdfb3be83cb761726afb33e0adaf32, with debug_info, not stripped
```