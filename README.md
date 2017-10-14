# Pancake Mahjong Server

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/mjpancake)
[![Build Status](https://travis-ci.org/rolevax/ih.svg?branch=develop)](https://travis-ci.org/rolevax/ih)

See the [client repository](https://github.com/rolevax/mjpancake)
for an introduction to Pancake Mahjong.

## Build

Requirement:

- Unix-like environment with common tools
- Go 1.8 or above
- [Proto Actor](https://github.com/AsynkronIT/protoactor-go)
- SWIG 3.0
- Other dependencies are go-gettable

Steps:

```
go get github.com/rolevax/ih
cd $GOPATH/src/github.com/rolevax/ih
make
```

## Run

```
docker-compose up -d
```

This will pull a pre-built image if their is no local build.


