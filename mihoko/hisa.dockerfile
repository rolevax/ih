FROM golang
COPY hisa /go/bin/hisa
CMD hisa -port 6171 -redis redis:6379 -db db:5432 -ryuuka hisa:6172 -toki toki:8900

