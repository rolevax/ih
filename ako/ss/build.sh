#!/bin/sh
protoc -I=. -I=${GOPATH}/src --gogoslick_out=plugins=grpc:. ss.proto 

