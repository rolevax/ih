# README

This repository contains submodules. 
Please don't forget to use `--recursive` when cloning. 

```
git clone --recursive https://bitbucket.org/rolevax/sakilogy-server.git
```
## Build

Requirement:

- Unix-like environment with common build tools (tested on Linux and macOS)
- Go
- SWIG 3.0
- All the other dependencies are go-gettable

Steps:

```
cd $GOPATH/src/bitbucket.org/rolevax/sakilogy-server
cd saki; make; cd ..
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
./sakilogy-server &
```

The log will output to stdout and stderr. 
Use `&>` to redirect if necessary.

