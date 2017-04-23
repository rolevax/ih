# Pancake Mahjong Server

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/mjpancake)

See the [client repository](https://github.com/mjpancake/mjpancake)
for an introduction to Pancake Mahjong.

## Build

Requirement:

- Unix-like environment with common build tools
- Go
- SWIG 3.0
- All the other dependencies are go-gettable

Steps:

```
go get github.com/mjpancake/hisa
cd $GOPATH/src/github.com/mjpancake/hisa
cd saki; make clean; make; cd ..
go build
```

## Run

Requirement:

- A running MySQL server with related tables already created
  (see the wiki for the database schema)
- A running Redis server with the default configurations
  (see `srv/redis.go`)

Then run the server:

```
./hisa &
```

The log will output to stdout and stderr. 
Use `&>` to redirect if necessary.

