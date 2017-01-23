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

Steps:

```
cd $GOPATH/src/bitbucket.org/rolevax/sakilogy-server
cd saki; make; cd ..
go build
```

## Run

Requires the correspending MySQL database running on the local machine.

```
./sakilogy-server
```

