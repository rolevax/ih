CXX=g++
CXXFLAGS=-std=c++11 -c -fPIC -O2
CXXSRC=../libsaki/*.cpp
BUILD_DIR=./build

AR=ar
ARFLAGS=crf

SED_REGEX='s/\*\//\#cgo LDFLAGS: \$${SRCDIR}\/build\/libsaki.a\n\*\//g'

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
	@sed -i ${SED_REGEX} -- saki.go
	@echo -e '\e[OK\r[DONE] saki.go saki_wrap.cxx'

${BUILD_DIR}:
	@mkdir -p ${BUILD_DIR}

clean:
	@rm -f saki.go saki_wrap.cxx
	@rm -rf ./build


