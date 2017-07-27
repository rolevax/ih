# 
# Makefile of Pancake Mahjong Server
#

HISA=hisa/hisa
TERU=teru/teru

all: hisa teru

.PHONY: saki
.PHONY: hisa
.PHONY: teru

hisa: ${HISA}

teru: ${TERU}

saki:
	${MAKE} -C saki

${HISA}: saki
	cd hisa;\
		go build;

${TERU}: saki
	pwd
	cd teru;\
		go build;

