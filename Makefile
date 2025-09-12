GOPATH = $(shell go env GOPATH)
VERSION = $(shell git describe --tags )
TODAY = $(shell date -u +'%Y-%m-%d')
build:
	go build -ldflags "-X main.version=${VERSION} -X main.date=${TODAY}" -o activitymonitor main.go

run: 
	go run main.go

clean:
	rm -f activitymonitor

install: build
	install -m 755 activitymonitor ${GOPATH}/bin/activitymonitor
