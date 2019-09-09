PROJECT="probe-host"

default:
	echo ${PROJECT}

install:
	mkdir -p bin
	export GO111MODULE=on
	export GOPROXY=https://goproxy.cn,direct
	CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/probe-host *.go

clean: 
	echo "clean probe-host binary"
	rm -rf bin/probe-host

.PHONY: default install clean

