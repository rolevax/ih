# 
# Makefile for docker images
#

all: teru-image hisa-image

teru-image:
	mkdir -p teru-tmp
	cp ../teru/teru ./teru-tmp/teru
	cp ./teru.dockerfile ./teru-tmp/Dockerfile
	docker build --rm -t rolevax/teru ./teru-tmp
	rm -r teru-tmp

hisa-image:
	mkdir -p hisa-tmp
	cp ../hisa/hisa ./hisa-tmp/hisa
	cp ./hisa.dockerfile ./hisa-tmp/Dockerfile
	docker build --rm -t rolevax/hisa ./hisa-tmp
	rm -r hisa-tmp

