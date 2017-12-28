# 
# Makefile of Pancake Mahjong Server
#

BIN_HISA=hisa/hisa
BIN_TERU=teru/teru
BIN_TOKI=toki/toki

all: hisa teru toki

.PHONY: saki
.PHONY: mihoko
.PHONY: hisa
.PHONY: teru
.PHONY: toki

hisa: ${BIN_HISA}

teru: ${BIN_TERU}

toki: ${BIN_TOKI}

saki:
	${MAKE} -C saki

mihoko:
	docker run --rm \
	           -v ${PWD}:/go/src/github.com/rolevax/ih \
	           -v /tmp/cache/rolevax/ih:/go/pkg/linux_amd64/github.com/rolevax/ih \
	           rolevax/ih-builder \
	           make
	${MAKE} -C mihoko

${BIN_HISA}: toki
	cd hisa;\
		go build;

${BIN_TERU}:
	cd teru;\
		go build;

${BIN_TOKI}: saki
	cd toki;\
		go build;

.PHONY: clean

clean:
	${MAKE} -C saki clean
	rm -f ${BIN_HISA} ${BIN_TERU} ${BIN_TOKI}


