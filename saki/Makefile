CXX=g++
CXXFLAGS=-std=c++11 -c -fPIC -O2
CXXSRC=../libsaki/*.cpp
BUILD_DIR=./build

AR=ar
ARFLAGS=cr

# inject the CGO preamble before the '*/' end-of-commend
CGO_INJ=cgo_preamble.inject
SED_MAGIC=sed -n -e '/*\//r ${CGO_INJ}'\
		  -e 1x -e '2,$$x' -e '2,$$p' -e '$$x' -e '$$p'\
		  saki.go > tmp\
		  && mv tmp saki.go

all: libsaki swig

libsaki: ${BUILD_DIR}
	@echo -n 'compiling libsaki...'
	@cd ${BUILD_DIR};\
		${CXX} ${CXXFLAGS} ${CXXSRC};\
		${AR} ${ARFLAGS} libsaki.a *.o;
	@echo -e '\e[OK\r[DONE] build/libsaki.a'

swig:
	@echo -n 'generating swig outputs...'
	@swig -go -cgo -c++ -intgosize 64 saki.i
	@${SED_MAGIC}
	@echo -e '\e[OK\r[DONE] saki.go saki_wrap.cxx'

${BUILD_DIR}:
	@mkdir -p ${BUILD_DIR}

clean:
	@rm -f saki.go saki_wrap.cxx
	@rm -rf ./build


