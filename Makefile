BUILD_ENV := CGO_ENABLED=1
BUILT=`date +%FT%T%z`
GOVERSION=`go version`
OSVERSION=`uname -a`
GITCOMMIT=`git rev-parse --short HEAD`
VERSION = 1.0.0.1

LDFLAGS=-ldflags "-v -s -w -X main.version=${VERSION} -X main.built=${BUILT} -X 'main.goversion=${GOVERSION}' -X 'main.osversion=${OSVERSION}' -X 'main.gitcommit=${GITCOMMIT}'"

TARGET_EXEC:=webappserver

.PHONY: all clean setup bl bo bw

all: clean setup bl bo bw

clean:
	rm -rf build

setup:
	mkdir -p build/linux
	mkdir -p build/osx
	mkdir -p build/windows

bl: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -a -v -o build/linux/${TARGET_EXEC}

bo: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -a -v -o build/osx/${TARGET_EXEC}

bw: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -a -v -o build/windows/${TARGET_EXEC}.exe
