# 
# Makefile of Pancake Mahjong Server
#

BIN_HISA=hisa/hisa
BIN_TERU=teru/teru

all: hisa teru

.PHONY: saki
.PHONY: hisa
.PHONY: teru

hisa: ${BIN_HISA}

teru: ${BIN_TERU}

saki:
	${MAKE} -C saki

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


