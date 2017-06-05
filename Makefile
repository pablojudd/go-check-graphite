NAME="check-graphite"

all: test build

test: lint
	export GOPATH=${PWD}/lib:${PWD}/build; \
	cd ./lib/src/CheckGraphite; \
	go test

lint: test-deps
	build/bin/golint -set_exit_status main.go
	build/bin/golint -set_exit_status ./lib/src/*

test-deps:
	GOPATH=${PWD}/build go get github.com/golang/lint/golint
	GOPATH=${PWD}/build go get github.com/pablojudd/go-graphite-getmetrics

deps:
	GOPATH=${PWD}/build go get github.com/jessevdk/go-flags
	GOPATH=${PWD}/build go get github.com/pablojudd/go-graphite-getmetrics

build: deps
	export GOPATH=${PWD}/lib:${PWD}/build; \
	go build -ldflags "-s -w" -o ${PWD}/${NAME} *.go

clean:
	rm -rf build
	for i in `/bin/ls ./lib` ; do if [ $$i != "src" ] ; then rm -rfv ./lib/$$i ; fi ; done
	rm -f ./${NAME}
	go clean

