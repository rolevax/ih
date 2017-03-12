# README

```
## Build

Requirement:

- Unix-like environment with common build tools (tested on Linux and macOS)
- Go
- SWIG 3.0
- All the other dependencies are go-gettable

Steps:

```
git clone --recursive https://github.com/mjpancake/mjpancake-server.git
cd $GOPATH/src/github.com/mjpancake/mjpancake-server
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
./mjpancake-server &
```

The log will output to stdout and stderr. 
Use `&>` to redirect if necessary.

