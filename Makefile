# 
# Makefile of Pancake Mahjong Server
#

BIN_HISA=hisa/hisa
BIN_TERU=teru/teru

all: hisa teru

.PHONY: saki
.PHONY: mihoko
.PHONY: hisa
.PHONY: teru

hisa: ${BIN_HISA}

teru: ${BIN_TERU}

saki:
	${MAKE} -C saki

mihoko:
	docker run --rm \
	           -v ${PWD}:/go/src/github.com/rolevax/ih \
	           -v /tmp/cache/rolevax/ih:/go/pkg/linux_amd64/github.com/rolevax/ih \
	           rolevax/ih-builder \
	           make teru
	${MAKE} -C mihoko

${BIN_HISA}: saki
	cd hisa;\
		go build;

${BIN_TERU}: saki
	cd teru;\
		go build;

.PHONY: clean

clean:
	${MAKE} -C saki clean
	rm -f ${BIN_HISA} ${BIN_TERU}


