FROM golang
COPY toki /go/bin/toki
CMD toki -addr toki:8900

