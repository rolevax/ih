FROM golang
COPY teru /go/bin/teru
CMD teru -redis redis:6379 -db db:5432

