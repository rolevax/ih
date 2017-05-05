# Pancake Mahjong Server

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/mjpancake)

See the [client repository](https://github.com/mjpancake/mjpancake)
for an introduction to Pancake Mahjong.

## Build

Requirement:

- Unix-like environment with common build tools
- Go 1.8 or above
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

- A running MySQL server with a well-built database.
  - To initialize a DB, start `mysql` as root and `source init.sql;`
- A running Redis server with the default configurations

Then run the server:

```
./hisa &
```

The log will output to stdout and stderr. 
Use `&>` to redirect if necessary.

