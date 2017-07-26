# Pancake Mahjong Server

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/mjpancake)

See the [client repository](https://github.com/mjpancake/mjpancake)
for an introduction to Pancake Mahjong.

## Build

Requirement:

- Unix-like environment with common build tools
- Go 1.8 or above
- [Proto Actor](https://github.com/AsynkronIT/protoactor-go)
- SWIG 3.0
- All the other dependencies are go-gettable

Steps:

```
go get github.com/mjpancake/ih
cd $GOPATH/src/github.com/mjpancake/ih
cd saki; make; cd ..
cd hisa; go install
```

## Run

Requirement:

- A running PostgreSQL server with required data
  - To import a dummy dataset, run `psql mako mako < mako/schema.pgsql`
- A running Redis server with the default configuration

Then run the server:

```
hisa &
```

The log will output to stdout and stderr. 
Use `&>` to redirect if necessary.

