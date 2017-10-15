# 
# Makefile for docker images
#

all: teru-image

teru-image:
	mkdir -p teru-tmp
	cp ../teru/teru ./teru-tmp/teru
	cp ./teru.dockerfile ./teru-tmp/Dockerfile
	docker build --rm -t rolevax/teru ./teru-tmp
	rm -r teru-tmp

