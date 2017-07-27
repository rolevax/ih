# 
# Makefile of ih/saki
#
# The default target does 3 things:
# 1. Compile libsaki.a
# 2. Generate SWIG code
# 3. Install ih/saki to GOPATH
#
# ih/saki builds very slow because it contains c++ code.
# We install it to GOPATH to fast-build other parts of ih.
#

BUILD_DIR=./build

CXX=g++
CXXFLAGS=-std=c++11 -c -fPIC -O2
CXXSRC=./libsaki/*.cpp
CXXHDR=./libsaki/*.h

AR=ar
ARFLAGS=cr

SWIG_INPUT=saki.i
SWIG_OUTPUT_GO=saki.go
SWIG_OUTPUT_CXX=saki_wrap.cxx

CGO_INJ=cgo_preamble.inject
CGO_SRC=*.cpp
CGO_HDR=*.h *.hpp

# inject the CGO preamble before the '*/' end-of-commend
SED_MAGIC=sed -n -e '/*\//r ${CGO_INJ}'\
		  -e 1x -e '2,$$x' -e '2,$$p' -e '$$x' -e '$$p'\
		  saki.go > tmp\
		  && mv tmp saki.go

GO_SAKI_PKG_NAME=github.com/mjpancake/ih/saki
GO_OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
GO_SAKI_PKG=${GOPATH}/pkg/${GO_OS_ARCH}/${GO_SAKI_PKG_NAME}.a

all: ${GO_SAKI_PKG}

${GO_SAKI_PKG}: ${BUILD_DIR}/libsaki.a\
  ${SWIG_OUTPUT_GO} ${SWIG_OUTPUT_CXX}\
  ${CGO_HDR} ${CGO_SRC}
	@echo -n 'install Go package ${GO_SAKI_PKG_NAME}... '
	@go install
	@echo 'done'

${BUILD_DIR}/libsaki.a: ${CXXHDR} ${CXXSRC} | ${BUILD_DIR}
	@echo -n 'compile ./build/libsaki.a... '
	@cd ${BUILD_DIR};\
		${CXX} ${CXXFLAGS} ../${CXXSRC};\
		${AR} ${ARFLAGS} libsaki.a *.o;
	@echo 'done'

${SWIG_OUTPUT_CXX}: ${SWIG_OUTPUT_GO}

${SWIG_OUTPUT_GO}: ${SWIG_INPUT} ${CGO_INJ}
	@echo -n 'generate swig code... '
	@swig -go -cgo -c++ -intgosize 64 ${SWIG_INPUT}
	@${SED_MAGIC}
	@echo 'done'

${BUILD_DIR}:
	@mkdir -p ${BUILD_DIR}

clean:
	@rm -f ${SWIG_OUTPUT_GO} ${SWIG_OUTPUT_CXX}
	@rm -rf ${BUILD_DIR}
	@rm -f ${GO_SAKI_PKG}


